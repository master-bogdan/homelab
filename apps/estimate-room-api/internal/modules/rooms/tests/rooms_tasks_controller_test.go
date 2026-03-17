package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	roomsdto "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/dto"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupRoomsTasksTest(t *testing.T) (*chi.Mux, *bun.DB) {
	t.Helper()

	db := testutils.SetupTestDB(t)

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
			invitations,
			votes,
			tasks,
			team_members,
			teams,
			room_participants,
			rooms,
			oauth2_access_tokens,
			oauth2_refresh_tokens,
			oauth2_auth_codes,
			oauth2_oidc_sessions,
			users,
			oauth2_clients
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	router := chi.NewRouter()
	authService := oauth2.NewAuthServiceFromDB(testutils.TestTokenKey, db)
	wsService := ws.NewService(nil, "test-room-events")

	router.Route("/api/v1", func(r chi.Router) {
		invitesModule := invites.NewInvitesModule(invites.InvitesModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
		})

		rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:         r,
			DB:             db,
			WsService:      wsService,
			AuthService:    authService,
			InvitesService: invitesModule.Service,
		})
	})

	return router, db
}

func seedRoom(t *testing.T, db *bun.DB, adminUserID string) string {
	t.Helper()

	roomID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO rooms (room_id, code, name, admin_user_id, deck)
		VALUES ($1, $2, $3, $4, '{"name":"Fibonacci","kind":"FIBONACCI","values":["0","1","2","3","5","8"]}'::jsonb)
	`, roomID, "code-"+roomID[:8], "Test Room", adminUserID)
	if err != nil {
		t.Fatalf("failed to insert room: %v", err)
	}

	_, err = db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, user_id, role)
		VALUES ($1, $2, $3, 'ADMIN')
	`, uuid.NewString(), roomID, adminUserID)
	if err != nil {
		t.Fatalf("failed to insert admin participant: %v", err)
	}

	return roomID
}

func createAccessToken(t *testing.T, db *bun.DB) (string, string) {
	t.Helper()

	return createAccessTokenForEmail(t, db, uuid.NewString()+"@example.com")
}

func createAccessTokenForEmail(t *testing.T, db *bun.DB, email string) (string, string) {
	t.Helper()

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, email, "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-rooms")

	svc := testutils.NewOauth2Service(db)
	tokens, err := svc.GenerateTokenPair(context.Background(), userID, clientID, sessionID, []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	return tokens.AccessToken, userID
}

func seedTeamForRoomTest(t *testing.T, db *bun.DB, ownerUserID, name string) string {
	t.Helper()

	teamID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO teams (team_id, name, owner_user_id)
		VALUES ($1, $2, $3)
	`, teamID, name, ownerUserID)
	if err != nil {
		t.Fatalf("failed to insert team: %v", err)
	}

	_, err = db.ExecContext(context.Background(), `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, 'OWNER')
	`, teamID, ownerUserID)
	if err != nil {
		t.Fatalf("failed to insert team owner: %v", err)
	}

	return teamID
}

func seedTeamMemberForRoomTest(t *testing.T, db *bun.DB, teamID, userID, role string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
	`, teamID, userID, role)
	if err != nil {
		t.Fatalf("failed to insert team member: %v", err)
	}
}

func seedUserWithoutEmailForRoomTest(t *testing.T, db *bun.DB) string {
	t.Helper()

	userID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO users (user_id, display_name)
		VALUES ($1, $2)
	`, userID, "No Email")
	if err != nil {
		t.Fatalf("failed to insert no-email user: %v", err)
	}

	return userID
}

func createRoomViaAPI(
	t *testing.T,
	router *chi.Mux,
	accessToken string,
	requestBody string,
) roomsdto.CreateRoomResponse {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/", bytes.NewReader([]byte(requestBody)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when creating room, got %d: %s", rr.Code, rr.Body.String())
	}

	var payload roomsdto.CreateRoomResponse
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode create room response: %v", err)
	}

	if payload.Room == nil || payload.Room.RoomID == "" {
		t.Fatal("expected room id in create room response")
	}

	return payload
}

func TestCreateRoom_DoesNotExposeRawServiceErrors(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, _ := createAccessToken(t, db)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/", bytes.NewReader([]byte(`{
		"name":"Invalid Deck Room",
		"deck":{"name":"","kind":"","values":[]}
	}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d: %s", rr.Code, rr.Body.String())
	}

	var httpErr apperrors.HttpError
	if err := json.NewDecoder(rr.Body).Decode(&httpErr); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if httpErr.Detail != "failed to create room" {
		t.Fatalf("expected sanitized detail, got %q", httpErr.Detail)
	}
	if httpErr.Detail == "invalid deck" {
		t.Fatal("expected raw service error to stay internal")
	}
}

func TestCreateRoom_DoesNotCreateShareLinkUnlessRequested(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, _ := createAccessToken(t, db)

	response := createRoomViaAPI(t, router, accessToken, `{"name":"No Share Link Room"}`)
	if response.InviteToken != "" {
		t.Fatalf("expected no inviteToken, got %q", response.InviteToken)
	}
	if response.ShareLink != nil {
		t.Fatalf("expected no share link, got %#v", response.ShareLink)
	}

	var invitationCount int
	err := db.NewSelect().
		Table("invitations").
		ColumnExpr("COUNT(*)").
		Where("room_id = ?", response.Room.RoomID).
		Scan(context.Background(), &invitationCount)
	if err != nil {
		t.Fatalf("failed to count invitations: %v", err)
	}
	if invitationCount != 0 {
		t.Fatalf("expected no invitations, got %d", invitationCount)
	}
}

func TestCreateRoom_FansOutTeamMembersAndExplicitEmails(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createAccessTokenForEmail(t, db, "owner@example.com")
	_, memberUserID := createAccessTokenForEmail(t, db, "member@example.com")
	noEmailUserID := seedUserWithoutEmailForRoomTest(t, db)
	teamID := seedTeamForRoomTest(t, db, ownerUserID, "Platform Team")
	seedTeamMemberForRoomTest(t, db, teamID, memberUserID, "MEMBER")
	seedTeamMemberForRoomTest(t, db, teamID, noEmailUserID, "MEMBER")

	response := createRoomViaAPI(t, router, ownerToken, `{
		"name":"Fanout Room",
		"inviteTeamId":"`+teamID+`",
		"inviteEmails":["member@example.com","owner@example.com","outside@example.com"],
		"createShareLink":true
	}`)

	if response.ShareLink == nil || response.InviteToken == "" {
		t.Fatalf("expected share link in create response, got %#v", response.ShareLink)
	}
	if len(response.EmailInvites) != 2 {
		t.Fatalf("expected 2 email invites, got %d", len(response.EmailInvites))
	}

	emailInvites := make(map[string]invitesdto.InvitationWithTokenResponse, len(response.EmailInvites))
	for _, invite := range response.EmailInvites {
		if invite.Kind != "ROOM_EMAIL" {
			t.Fatalf("expected ROOM_EMAIL invite kind, got %s", invite.Kind)
		}
		if invite.RoomID == nil || *invite.RoomID != response.Room.RoomID {
			t.Fatalf("expected room id %s, got %#v", response.Room.RoomID, invite.RoomID)
		}
		if invite.InvitedEmail == nil {
			t.Fatalf("expected invited email, got %#v", invite)
		}
		emailInvites[*invite.InvitedEmail] = invite
	}

	if _, exists := emailInvites["member@example.com"]; !exists {
		t.Fatalf("expected member@example.com invite, got %#v", emailInvites)
	}
	if _, exists := emailInvites["outside@example.com"]; !exists {
		t.Fatalf("expected outside@example.com invite, got %#v", emailInvites)
	}
	if _, exists := emailInvites["owner@example.com"]; exists {
		t.Fatalf("did not expect self invite, got %#v", emailInvites)
	}

	if len(response.SkippedRecipients) != 2 {
		t.Fatalf("expected 2 skipped recipients, got %#v", response.SkippedRecipients)
	}

	reasonsByUserID := make(map[string]string, len(response.SkippedRecipients))
	reasonsByEmail := make(map[string]string, len(response.SkippedRecipients))
	for _, skipped := range response.SkippedRecipients {
		if skipped.UserID != nil {
			reasonsByUserID[*skipped.UserID] = skipped.Reason
		}
		if skipped.Email != nil {
			reasonsByEmail[*skipped.Email] = skipped.Reason
		}
	}

	if reasonsByUserID[noEmailUserID] != "missing_email" {
		t.Fatalf("expected missing_email skip for no-email user, got %#v", response.SkippedRecipients)
	}
	if reasonsByEmail["owner@example.com"] != "self" {
		t.Fatalf("expected self skip for owner email, got %#v", response.SkippedRecipients)
	}

	var roomEmailInviteCount int
	err := db.NewSelect().
		Table("invitations").
		ColumnExpr("COUNT(*)").
		Where("room_id = ?", response.Room.RoomID).
		Where("kind = 'ROOM_EMAIL'").
		Scan(context.Background(), &roomEmailInviteCount)
	if err != nil {
		t.Fatalf("failed to count room email invites: %v", err)
	}
	if roomEmailInviteCount != 2 {
		t.Fatalf("expected 2 persisted room email invites, got %d", roomEmailInviteCount)
	}
}

func TestCreateRoom_RejectsInviteTeamIdForNonOwner(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	_, ownerUserID := createAccessTokenForEmail(t, db, "owner@example.com")
	memberToken, memberUserID := createAccessTokenForEmail(t, db, "member@example.com")
	teamID := seedTeamForRoomTest(t, db, ownerUserID, "Restricted Team")
	seedTeamMemberForRoomTest(t, db, teamID, memberUserID, "MEMBER")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/", bytes.NewReader([]byte(`{
		"name":"Forbidden Fanout",
		"inviteTeamId":"`+teamID+`"
	}`)))
	req.Header.Set("Authorization", "Bearer "+memberToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}

	var roomCount int
	err := db.NewSelect().
		Table("rooms").
		ColumnExpr("COUNT(*)").
		Scan(context.Background(), &roomCount)
	if err != nil {
		t.Fatalf("failed to count rooms: %v", err)
	}
	if roomCount != 0 {
		t.Fatalf("expected no room to be created, found %d", roomCount)
	}
}

func TestTasksCRUD(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	createReqBody := []byte(`{"title":"Initial task","description":"Estimate login","externalKey":"JIRA-101"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/"+roomID+"/tasks/", bytes.NewReader(createReqBody))
	createReq.Header.Set("Authorization", "Bearer "+accessToken)

	createRR := httptest.NewRecorder()
	router.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on create, got %d: %s", createRR.Code, createRR.Body.String())
	}

	var createdTask roomsmodels.RoomTaskModel
	if err := json.NewDecoder(createRR.Body).Decode(&createdTask); err != nil {
		t.Fatalf("failed to decode created task: %v", err)
	}

	if createdTask.TaskID == "" {
		t.Fatal("expected created task id")
	}
	if createdTask.RoomID != roomID {
		t.Fatalf("expected room id %s, got %s", roomID, createdTask.RoomID)
	}
	if createdTask.Status != "PENDING" {
		t.Fatalf("expected status PENDING, got %s", createdTask.Status)
	}
	if createdTask.IsActive {
		t.Fatal("expected created task to be inactive")
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+roomID+"/tasks/", nil)
	listReq.Header.Set("Authorization", "Bearer "+accessToken)

	listRR := httptest.NewRecorder()
	router.ServeHTTP(listRR, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list, got %d: %s", listRR.Code, listRR.Body.String())
	}

	var tasks []*roomsmodels.RoomTaskModel
	if err := json.NewDecoder(listRR.Body).Decode(&tasks); err != nil {
		t.Fatalf("failed to decode tasks list: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, nil)
	getReq.Header.Set("Authorization", "Bearer "+accessToken)

	getRR := httptest.NewRecorder()
	router.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on get, got %d: %s", getRR.Code, getRR.Body.String())
	}

	updateReqBody := []byte(`{"title":"Updated task"}`)
	updateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, bytes.NewReader(updateReqBody))
	updateReq.Header.Set("Authorization", "Bearer "+accessToken)

	updateRR := httptest.NewRecorder()
	router.ServeHTTP(updateRR, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on update, got %d: %s", updateRR.Code, updateRR.Body.String())
	}

	var updatedTask roomsmodels.RoomTaskModel
	if err := json.NewDecoder(updateRR.Body).Decode(&updatedTask); err != nil {
		t.Fatalf("failed to decode updated task: %v", err)
	}
	if updatedTask.Title != "Updated task" {
		t.Fatalf("expected updated title, got %s", updatedTask.Title)
	}
	if updatedTask.Status != "PENDING" {
		t.Fatalf("expected updated status PENDING, got %s", updatedTask.Status)
	}
	if updatedTask.IsActive {
		t.Fatal("expected metadata update to keep task inactive")
	}

	activateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, bytes.NewReader([]byte(`{"isActive":true}`)))
	activateReq.Header.Set("Authorization", "Bearer "+accessToken)

	activateRR := httptest.NewRecorder()
	router.ServeHTTP(activateRR, activateReq)

	if activateRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on activation, got %d: %s", activateRR.Code, activateRR.Body.String())
	}

	var activatedTask roomsmodels.RoomTaskModel
	if err := json.NewDecoder(activateRR.Body).Decode(&activatedTask); err != nil {
		t.Fatalf("failed to decode activated task: %v", err)
	}
	if activatedTask.Status != "VOTING" {
		t.Fatalf("expected activated status VOTING, got %s", activatedTask.Status)
	}
	if !activatedTask.IsActive {
		t.Fatal("expected activated task to be active")
	}

	_, err := db.ExecContext(context.Background(), `
			UPDATE task_rounds
			SET status = 'REVEALED', updated_at = NOW()
			WHERE task_id = $1 AND round_number = 1
		`, createdTask.TaskID)
	if err != nil {
		t.Fatalf("failed to mark round revealed: %v", err)
	}

	finalizeReq := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, bytes.NewReader([]byte(`{"finalEstimateValue":"5"}`)))
	finalizeReq.Header.Set("Authorization", "Bearer "+accessToken)

	finalizeRR := httptest.NewRecorder()
	router.ServeHTTP(finalizeRR, finalizeReq)

	if finalizeRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on finalize, got %d: %s", finalizeRR.Code, finalizeRR.Body.String())
	}

	var finalizedTask roomsmodels.RoomTaskModel
	if err := json.NewDecoder(finalizeRR.Body).Decode(&finalizedTask); err != nil {
		t.Fatalf("failed to decode finalized task: %v", err)
	}
	if finalizedTask.Status != "ESTIMATED" {
		t.Fatalf("expected finalized status ESTIMATED, got %s", finalizedTask.Status)
	}
	if finalizedTask.IsActive {
		t.Fatal("expected finalized task to be inactive")
	}
	if finalizedTask.FinalEstimateValue == nil || *finalizedTask.FinalEstimateValue != "5" {
		t.Fatalf("expected final estimate value 5, got %#v", finalizedTask.FinalEstimateValue)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+accessToken)

	deleteRR := httptest.NewRecorder()
	router.ServeHTTP(deleteRR, deleteReq)

	if deleteRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on delete, got %d: %s", deleteRR.Code, deleteRR.Body.String())
	}

	getMissingReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, nil)
	getMissingReq.Header.Set("Authorization", "Bearer "+accessToken)

	getMissingRR := httptest.NewRecorder()
	router.ServeHTTP(getMissingRR, getMissingReq)

	if getMissingRR.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found after delete, got %d: %s", getMissingRR.Code, getMissingRR.Body.String())
	}

	var httpErr apperrors.HttpError
	if err := json.NewDecoder(getMissingRR.Body).Decode(&httpErr); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if httpErr.Status != http.StatusNotFound {
		t.Fatalf("expected error status 404, got %d", httpErr.Status)
	}
}

func TestTasksUpdate_ActivatingTaskClearsPreviousActiveTask(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	createTask := func(title string) roomsmodels.RoomTaskModel {
		t.Helper()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/"+roomID+"/tasks/", bytes.NewReader([]byte(`{"title":"`+title+`"}`)))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on create, got %d: %s", rr.Code, rr.Body.String())
		}

		var task roomsmodels.RoomTaskModel
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("failed to decode task: %v", err)
		}

		return task
	}

	firstTask := createTask("Task One")
	secondTask := createTask("Task Two")

	activate := func(taskID string) roomsmodels.RoomTaskModel {
		t.Helper()

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID+"/tasks/"+taskID, bytes.NewReader([]byte(`{"isActive":true}`)))
		req.Header.Set("Authorization", "Bearer "+accessToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on activation, got %d: %s", rr.Code, rr.Body.String())
		}

		var task roomsmodels.RoomTaskModel
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("failed to decode activated task: %v", err)
		}

		return task
	}

	activatedFirst := activate(firstTask.TaskID)
	if !activatedFirst.IsActive || activatedFirst.Status != "VOTING" {
		t.Fatalf("expected first task to be active and VOTING, got active=%v status=%s", activatedFirst.IsActive, activatedFirst.Status)
	}

	activatedSecond := activate(secondTask.TaskID)
	if !activatedSecond.IsActive || activatedSecond.Status != "VOTING" {
		t.Fatalf("expected second task to be active and VOTING, got active=%v status=%s", activatedSecond.IsActive, activatedSecond.Status)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+roomID+"/tasks/", nil)
	listReq.Header.Set("Authorization", "Bearer "+accessToken)

	listRR := httptest.NewRecorder()
	router.ServeHTTP(listRR, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on list, got %d: %s", listRR.Code, listRR.Body.String())
	}

	var tasks []*roomsmodels.RoomTaskModel
	if err := json.NewDecoder(listRR.Body).Decode(&tasks); err != nil {
		t.Fatalf("failed to decode tasks list: %v", err)
	}

	taskByID := make(map[string]*roomsmodels.RoomTaskModel, len(tasks))
	for _, task := range tasks {
		taskByID[task.TaskID] = task
	}

	if taskByID[firstTask.TaskID] == nil || taskByID[secondTask.TaskID] == nil {
		t.Fatalf("expected both tasks in list, got %#v", taskByID)
	}
	if taskByID[firstTask.TaskID].IsActive || taskByID[firstTask.TaskID].Status != "PENDING" {
		t.Fatalf("expected first task reset to inactive PENDING, got active=%v status=%s", taskByID[firstTask.TaskID].IsActive, taskByID[firstTask.TaskID].Status)
	}
	if !taskByID[secondTask.TaskID].IsActive || taskByID[secondTask.TaskID].Status != "VOTING" {
		t.Fatalf("expected second task active VOTING, got active=%v status=%s", taskByID[secondTask.TaskID].IsActive, taskByID[secondTask.TaskID].Status)
	}
}

func TestTasksUpdate_TouchesRoomActivity(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/"+roomID+"/tasks/", bytes.NewReader([]byte(`{"title":"Activity task"}`)))
	createReq.Header.Set("Authorization", "Bearer "+accessToken)

	createRR := httptest.NewRecorder()
	router.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on create, got %d: %s", createRR.Code, createRR.Body.String())
	}

	var createdTask roomsmodels.RoomTaskModel
	if err := json.NewDecoder(createRR.Body).Decode(&createdTask); err != nil {
		t.Fatalf("failed to decode created task: %v", err)
	}

	oldActivityAt := time.Now().Add(-1 * time.Hour).UTC()
	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET last_activity_at = $2
		WHERE room_id = $1
	`, roomID, oldActivityAt); err != nil {
		t.Fatalf("failed to set room activity time: %v", err)
	}

	updateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID+"/tasks/"+createdTask.TaskID, bytes.NewReader([]byte(`{"title":"Activity task updated"}`)))
	updateReq.Header.Set("Authorization", "Bearer "+accessToken)

	updateRR := httptest.NewRecorder()
	router.ServeHTTP(updateRR, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on update, got %d: %s", updateRR.Code, updateRR.Body.String())
	}

	room := new(roomsmodels.RoomsModel)
	if err := db.NewSelect().
		Model(room).
		Where("r.room_id = ?", roomID).
		Limit(1).
		Scan(context.Background()); err != nil {
		t.Fatalf("failed to load room: %v", err)
	}

	if !room.LastActivityAt.After(oldActivityAt) {
		t.Fatalf("expected lastActivityAt to advance beyond %s, got %s", oldActivityAt, room.LastActivityAt)
	}
}

func TestTasksCRUD_RequiresRoomAdmin(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	_, ownerUserID := createAccessToken(t, db)
	roomID := seedRoom(t, db, ownerUserID)

	memberToken, _ := createAccessToken(t, db)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/rooms/"+roomID+"/tasks/", bytes.NewReader([]byte(`{"title":"Initial task"}`)))
	req.Header.Set("Authorization", "Bearer "+memberToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for non-admin, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestInviteFlow_AuthenticatedUserJoinsAsMember(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	adminToken, _ := createAccessToken(t, db)
	createResponse := createRoomViaAPI(t, router, adminToken, `{"name":"Invite Room","createShareLink":true}`)
	roomID := createResponse.Room.RoomID
	inviteToken := createResponse.InviteToken

	memberToken, memberUserID := createAccessToken(t, db)
	joinReq := httptest.NewRequest(http.MethodPost, "/api/v1/invites/"+inviteToken+"/accept", nil)
	joinReq.Header.Set("Authorization", "Bearer "+memberToken)

	joinRR := httptest.NewRecorder()
	router.ServeHTTP(joinRR, joinReq)

	if joinRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on member invite join, got %d: %s", joinRR.Code, joinRR.Body.String())
	}

	var joined struct {
		Participant roomsmodels.RoomParticipantModel `json:"participant"`
	}
	if err := json.NewDecoder(joinRR.Body).Decode(&joined); err != nil {
		t.Fatalf("failed to decode invite join response: %v", err)
	}
	if joined.Participant.UserID == nil || *joined.Participant.UserID != memberUserID {
		t.Fatalf("expected participant user id %s, got %#v", memberUserID, joined.Participant.UserID)
	}
	if joined.Participant.Role != roomsmodels.RoomParticipantRoleMember {
		t.Fatalf("expected MEMBER role, got %s", joined.Participant.Role)
	}

	getRoomReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+roomID, nil)
	getRoomReq.Header.Set("Authorization", "Bearer "+memberToken)
	getRoomRR := httptest.NewRecorder()
	router.ServeHTTP(getRoomRR, getRoomReq)

	if getRoomRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when joined member reads room, got %d: %s", getRoomRR.Code, getRoomRR.Body.String())
	}
}

func TestInviteFlow_GuestJoinsAndCanReadRoom(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	adminToken, _ := createAccessToken(t, db)
	createResponse := createRoomViaAPI(t, router, adminToken, `{"name":"Invite Room","createShareLink":true}`)
	roomID := createResponse.Room.RoomID
	inviteToken := createResponse.InviteToken

	joinReq := httptest.NewRequest(http.MethodPost, "/api/v1/invites/"+inviteToken+"/accept", bytes.NewReader([]byte(`{"guestName":"Guest One"}`)))
	joinRR := httptest.NewRecorder()
	router.ServeHTTP(joinRR, joinReq)

	if joinRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on guest join, got %d: %s", joinRR.Code, joinRR.Body.String())
	}

	var joined struct {
		Participant roomsmodels.RoomParticipantModel `json:"participant"`
	}
	if err := json.NewDecoder(joinRR.Body).Decode(&joined); err != nil {
		t.Fatalf("failed to decode guest join response: %v", err)
	}

	if joined.Participant.GuestName == nil || *joined.Participant.GuestName != "Guest One" {
		t.Fatalf("expected guest name to be set, got %#v", joined.Participant.GuestName)
	}
	if joined.Participant.Role != roomsmodels.RoomParticipantRoleGuest {
		t.Fatalf("expected GUEST role, got %s", joined.Participant.Role)
	}

	cookies := joinRR.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected guest join to set a cookie")
	}

	getRoomReq := httptest.NewRequest(http.MethodGet, "/api/v1/rooms/"+roomID, nil)
	getRoomReq.AddCookie(cookies[0])
	getRoomRR := httptest.NewRecorder()
	router.ServeHTTP(getRoomRR, getRoomReq)

	if getRoomRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when guest reads room, got %d: %s", getRoomRR.Code, getRoomRR.Body.String())
	}
}
