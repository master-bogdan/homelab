package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/history"
	historydto "github.com/master-bogdan/estimate-room-api/internal/modules/history/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupHistoryTest(t *testing.T) (*chi.Mux, *bun.DB) {
	t.Helper()

	db := testutils.SetupTestDB(t)

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
			user_session_rewards,
			invitations,
			votes,
			tasks,
			task_rounds,
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
		t.Fatalf("failed to truncate history test tables: %v", err)
	}

	router := chi.NewRouter()
	authService := oauth2.NewOauth2SessionAuthServiceFromDB(testutils.TestTokenKey, db)

	router.Route("/api/v1", func(r chi.Router) {
		history.NewHistoryModule(history.HistoryModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
		})
	})

	return router, db
}

func createHistoryAccessToken(t *testing.T, db *bun.DB, email string) (string, string) {
	t.Helper()

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, email, "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-history")

	svc := testutils.NewOauth2Service(db)
	tokens, err := svc.GenerateTokenPair(context.Background(), userID, clientID, sessionID, []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	return tokens.AccessToken, userID
}

func seedHistoryTeam(t *testing.T, db *bun.DB, ownerUserID, name string) string {
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

func seedHistoryTeamMember(t *testing.T, db *bun.DB, teamID, userID, role string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
	`, teamID, userID, role)
	if err != nil {
		t.Fatalf("failed to insert team member: %v", err)
	}
}

func seedHistoryRoom(
	t *testing.T,
	db *bun.DB,
	roomID, name, adminUserID string,
	teamID *string,
	status string,
	createdAt, lastActivityAt time.Time,
	finishedAt *time.Time,
) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO rooms (room_id, code, name, admin_user_id, team_id, deck, status, created_at, last_activity_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, '{"name":"Fibonacci","kind":"FIBONACCI","values":["1","2","3","5","8"]}'::jsonb, $6, $7, $8, $9)
	`, roomID, "code-"+roomID[:8], name, adminUserID, teamID, status, createdAt.UTC(), lastActivityAt.UTC(), finishedAt)
	if err != nil {
		t.Fatalf("failed to insert room %s: %v", roomID, err)
	}
}

func seedHistoryParticipant(
	t *testing.T,
	db *bun.DB,
	roomID, userID string,
	role string,
	joinedAt time.Time,
) {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`, participantID, roomID, userID, role, joinedAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert participant: %v", err)
	}
}

func seedHistoryParticipantWithID(
	t *testing.T,
	db *bun.DB,
	roomID, userID string,
	role string,
	joinedAt time.Time,
) string {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`, participantID, roomID, userID, role, joinedAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert participant: %v", err)
	}

	return participantID
}

func seedHistoryGuestParticipant(
	t *testing.T,
	db *bun.DB,
	roomID, guestName string,
	role string,
	joinedAt time.Time,
) {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, guest_name, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`, participantID, roomID, guestName, role, joinedAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert guest participant: %v", err)
	}
}

func seedHistoryGuestParticipantWithID(
	t *testing.T,
	db *bun.DB,
	roomID, guestName string,
	role string,
	joinedAt time.Time,
) string {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, guest_name, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`, participantID, roomID, guestName, role, joinedAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert guest participant: %v", err)
	}

	return participantID
}

func seedHistoryTask(t *testing.T, db *bun.DB, roomID, title, status string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO tasks (task_id, room_id, title, status, is_active)
		VALUES ($1, $2, $3, $4, false)
	`, uuid.NewString(), roomID, title, status)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}
}

func seedHistoryTaskWithID(
	t *testing.T,
	db *bun.DB,
	roomID, title, status string,
	isActive bool,
	finalEstimateValue *string,
	createdAt, updatedAt time.Time,
) string {
	t.Helper()

	taskID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO tasks (task_id, room_id, title, status, is_active, final_estimate_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, taskID, roomID, title, status, isActive, finalEstimateValue, createdAt.UTC(), updatedAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert task with id: %v", err)
	}

	return taskID
}

func seedHistoryTaskRound(
	t *testing.T,
	db *bun.DB,
	taskID string,
	roundNumber int,
	status string,
	eligibleParticipantIDs []string,
	createdAt, updatedAt time.Time,
) {
	t.Helper()

	jsonEligible := `[]`
	if len(eligibleParticipantIDs) > 0 {
		jsonEligible = `["` + strings.Join(eligibleParticipantIDs, `","`) + `"]`
	}

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO task_rounds (task_id, round_number, eligible_participant_ids, status, created_at, updated_at)
		VALUES ($1, $2, $3::jsonb, $4, $5, $6)
	`, taskID, roundNumber, jsonEligible, status, createdAt.UTC(), updatedAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert task round: %v", err)
	}
}

func seedHistoryVote(
	t *testing.T,
	db *bun.DB,
	taskID, participantID string,
	roundNumber int,
	value string,
	createdAt time.Time,
) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO votes (votes_id, task_id, participant_id, value, round_number, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.NewString(), taskID, participantID, value, roundNumber, createdAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert vote: %v", err)
	}
}

func TestListMySessions_ReturnsSessionsWithCountsAndOrdering(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	token, userID := createHistoryAccessToken(t, db, "history-user@example.com")
	_, otherAdminUserID := createHistoryAccessToken(t, db, "other-admin@example.com")
	_, outsiderUserID := createHistoryAccessToken(t, db, "outsider@example.com")

	teamID := seedHistoryTeam(t, db, userID, "Platform")

	activeAdminRoomID := uuid.NewString()
	activeAdminCreatedAt := time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC)
	activeAdminActivityAt := time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC)
	seedHistoryRoom(t, db, activeAdminRoomID, "Admin Active Room", userID, &teamID, "ACTIVE", activeAdminCreatedAt, activeAdminActivityAt, nil)
	seedHistoryParticipant(t, db, activeAdminRoomID, userID, "ADMIN", activeAdminCreatedAt)
	seedHistoryGuestParticipant(t, db, activeAdminRoomID, "Guest Reviewer", "GUEST", activeAdminCreatedAt.Add(5*time.Minute))
	seedHistoryTask(t, db, activeAdminRoomID, "Backend task", "ESTIMATED")
	seedHistoryTask(t, db, activeAdminRoomID, "Frontend task", "PENDING")

	finishedParticipantRoomID := uuid.NewString()
	finishedCreatedAt := time.Date(2026, 3, 16, 8, 0, 0, 0, time.UTC)
	finishedActivityAt := time.Date(2026, 3, 16, 8, 45, 0, 0, time.UTC)
	finishedAt := time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC)
	seedHistoryRoom(t, db, finishedParticipantRoomID, "Finished Session", otherAdminUserID, nil, "FINISHED", finishedCreatedAt, finishedActivityAt, &finishedAt)
	seedHistoryParticipant(t, db, finishedParticipantRoomID, otherAdminUserID, "ADMIN", finishedCreatedAt)
	seedHistoryParticipant(t, db, finishedParticipantRoomID, userID, "MEMBER", finishedCreatedAt.Add(2*time.Minute))
	seedHistoryTask(t, db, finishedParticipantRoomID, "API task", "ESTIMATED")

	activeParticipantRoomID := uuid.NewString()
	activeParticipantCreatedAt := time.Date(2026, 3, 17, 8, 0, 0, 0, time.UTC)
	activeParticipantActivityAt := time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC)
	seedHistoryRoom(t, db, activeParticipantRoomID, "Latest Active Session", otherAdminUserID, nil, "ACTIVE", activeParticipantCreatedAt, activeParticipantActivityAt, nil)
	seedHistoryParticipant(t, db, activeParticipantRoomID, otherAdminUserID, "ADMIN", activeParticipantCreatedAt)
	seedHistoryParticipant(t, db, activeParticipantRoomID, userID, "MEMBER", activeParticipantCreatedAt.Add(1*time.Minute))

	unrelatedRoomID := uuid.NewString()
	unrelatedCreatedAt := time.Date(2026, 3, 14, 12, 0, 0, 0, time.UTC)
	unrelatedActivityAt := time.Date(2026, 3, 14, 13, 0, 0, 0, time.UTC)
	seedHistoryRoom(t, db, unrelatedRoomID, "Unrelated Room", outsiderUserID, nil, "ACTIVE", unrelatedCreatedAt, unrelatedActivityAt, nil)
	seedHistoryParticipant(t, db, unrelatedRoomID, outsiderUserID, "ADMIN", unrelatedCreatedAt)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/me/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.PaginatedResponse[historydto.SessionListItem]
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Total != 3 {
		t.Fatalf("expected total 3, got %d", response.Total)
	}
	if len(response.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(response.Items))
	}
	if response.Page != 1 || response.PageSize != 20 {
		t.Fatalf("expected default pagination 1/20, got %d/%d", response.Page, response.PageSize)
	}

	if response.Items[0].RoomID != activeParticipantRoomID {
		t.Fatalf("expected latest active participant room first, got %s", response.Items[0].RoomID)
	}
	if response.Items[0].Role != "PARTICIPANT" {
		t.Fatalf("expected PARTICIPANT role, got %s", response.Items[0].Role)
	}

	if response.Items[1].RoomID != finishedParticipantRoomID {
		t.Fatalf("expected finished participant room second, got %s", response.Items[1].RoomID)
	}
	if response.Items[1].Status != "FINISHED" {
		t.Fatalf("expected FINISHED status, got %s", response.Items[1].Status)
	}
	if response.Items[1].ApproxDurationSeconds != int64(finishedAt.Sub(finishedCreatedAt).Seconds()) {
		t.Fatalf("expected finished duration %d, got %d", int64(finishedAt.Sub(finishedCreatedAt).Seconds()), response.Items[1].ApproxDurationSeconds)
	}

	adminItem := response.Items[2]
	if adminItem.RoomID != activeAdminRoomID {
		t.Fatalf("expected admin room third, got %s", adminItem.RoomID)
	}
	if adminItem.Role != "ADMIN" {
		t.Fatalf("expected ADMIN role, got %s", adminItem.Role)
	}
	if adminItem.TeamID == nil || *adminItem.TeamID != teamID {
		t.Fatalf("expected team id %s, got %#v", teamID, adminItem.TeamID)
	}
	if adminItem.ParticipantsCount != 2 {
		t.Fatalf("expected 2 participants, got %d", adminItem.ParticipantsCount)
	}
	if adminItem.TasksCount != 2 {
		t.Fatalf("expected 2 tasks, got %d", adminItem.TasksCount)
	}
	if adminItem.EstimatedTasksCount != 1 {
		t.Fatalf("expected 1 estimated task, got %d", adminItem.EstimatedTasksCount)
	}
	if adminItem.ApproxDurationSeconds != int64(activeAdminActivityAt.Sub(activeAdminCreatedAt).Seconds()) {
		t.Fatalf("expected active duration %d, got %d", int64(activeAdminActivityAt.Sub(activeAdminCreatedAt).Seconds()), adminItem.ApproxDurationSeconds)
	}
}

func TestListMySessions_FiltersByStatusAndRole(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	token, userID := createHistoryAccessToken(t, db, "history-user@example.com")
	_, otherAdminUserID := createHistoryAccessToken(t, db, "other-admin@example.com")

	finishedParticipantRoomID := uuid.NewString()
	finishedCreatedAt := time.Date(2026, 3, 16, 8, 0, 0, 0, time.UTC)
	finishedActivityAt := time.Date(2026, 3, 16, 8, 45, 0, 0, time.UTC)
	finishedAt := time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC)
	seedHistoryRoom(t, db, finishedParticipantRoomID, "Finished Session", otherAdminUserID, nil, "FINISHED", finishedCreatedAt, finishedActivityAt, &finishedAt)
	seedHistoryParticipant(t, db, finishedParticipantRoomID, otherAdminUserID, "ADMIN", finishedCreatedAt)
	seedHistoryParticipant(t, db, finishedParticipantRoomID, userID, "MEMBER", finishedCreatedAt.Add(time.Minute))

	activeAdminRoomID := uuid.NewString()
	activeAdminCreatedAt := time.Date(2026, 3, 17, 8, 0, 0, 0, time.UTC)
	activeAdminActivityAt := time.Date(2026, 3, 17, 10, 0, 0, 0, time.UTC)
	seedHistoryRoom(t, db, activeAdminRoomID, "Admin Active", userID, nil, "ACTIVE", activeAdminCreatedAt, activeAdminActivityAt, nil)
	seedHistoryParticipant(t, db, activeAdminRoomID, userID, "ADMIN", activeAdminCreatedAt)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/me/sessions?status=finished&role=participant", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.PaginatedResponse[historydto.SessionListItem]
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Total != 1 || len(response.Items) != 1 {
		t.Fatalf("expected 1 filtered item, got total=%d len=%d", response.Total, len(response.Items))
	}
	if response.Items[0].RoomID != finishedParticipantRoomID {
		t.Fatalf("expected finished participant room %s, got %s", finishedParticipantRoomID, response.Items[0].RoomID)
	}
	if response.Items[0].Role != "PARTICIPANT" {
		t.Fatalf("expected PARTICIPANT role, got %s", response.Items[0].Role)
	}
}

func TestListMySessions_AllStatusFilterReturnsOK(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	token, userID := createHistoryAccessToken(t, db, "history-all-status@example.com")

	roomID := uuid.NewString()
	createdAt := time.Date(2026, 3, 18, 8, 0, 0, 0, time.UTC)
	lastActivityAt := time.Date(2026, 3, 18, 9, 0, 0, 0, time.UTC)
	seedHistoryRoom(t, db, roomID, "All Status Room", userID, nil, "ACTIVE", createdAt, lastActivityAt, nil)
	seedHistoryParticipant(t, db, roomID, userID, "ADMIN", createdAt)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/me/sessions?status=ALL", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for status=ALL, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.PaginatedResponse[historydto.SessionListItem]
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Total != 1 || len(response.Items) != 1 {
		t.Fatalf("expected one session for status=ALL, got total=%d len=%d", response.Total, len(response.Items))
	}
}

func TestListMySessions_PaginatesResults(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	token, userID := createHistoryAccessToken(t, db, "history-user@example.com")

	for i := 0; i < 3; i++ {
		roomID := uuid.NewString()
		createdAt := time.Date(2026, 3, 17-i, 8, 0, 0, 0, time.UTC)
		lastActivityAt := createdAt.Add(30 * time.Minute)
		seedHistoryRoom(t, db, roomID, "Room", userID, nil, "ACTIVE", createdAt, lastActivityAt, nil)
		seedHistoryParticipant(t, db, roomID, userID, "ADMIN", createdAt)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/me/sessions?page=2&pageSize=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.PaginatedResponse[historydto.SessionListItem]
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Total != 3 {
		t.Fatalf("expected total 3, got %d", response.Total)
	}
	if response.Page != 2 || response.PageSize != 1 {
		t.Fatalf("expected pagination 2/1, got %d/%d", response.Page, response.PageSize)
	}
	if len(response.Items) != 1 {
		t.Fatalf("expected 1 item on page 2, got %d", len(response.Items))
	}
}

func TestListMySessions_InvalidFiltersReturnBadRequest(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	token, _ := createHistoryAccessToken(t, db, "history-user@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/me/sessions?status=weird&role=nope", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}

	var response apperrors.HttpError
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if response.Status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Status)
	}
}

func TestListTeamSessions_RequiresTeamOwnerAndSupportsStatusFilter(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createHistoryAccessToken(t, db, "owner@example.com")
	memberToken, memberUserID := createHistoryAccessToken(t, db, "member@example.com")
	outsiderToken, _ := createHistoryAccessToken(t, db, "outsider@example.com")

	teamID := seedHistoryTeam(t, db, ownerUserID, "Platform")
	seedHistoryTeamMember(t, db, teamID, memberUserID, "MEMBER")

	finishedAt := time.Date(2026, 3, 16, 11, 0, 0, 0, time.UTC)
	finishedRoomID := uuid.NewString()
	seedHistoryRoom(
		t,
		db,
		finishedRoomID,
		"Finished Team Room",
		memberUserID,
		&teamID,
		"FINISHED",
		time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC),
		&finishedAt,
	)
	seedHistoryParticipant(t, db, finishedRoomID, memberUserID, "ADMIN", time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC))

	activeRoomID := uuid.NewString()
	seedHistoryRoom(
		t,
		db,
		activeRoomID,
		"Active Team Room",
		ownerUserID,
		&teamID,
		"ACTIVE",
		time.Date(2026, 3, 17, 8, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 17, 9, 30, 0, 0, time.UTC),
		nil,
	)
	seedHistoryParticipant(t, db, activeRoomID, ownerUserID, "ADMIN", time.Date(2026, 3, 17, 8, 0, 0, 0, time.UTC))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/teams/"+teamID+"/sessions?status=finished", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.PaginatedResponse[historydto.SessionListItem]
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Total != 1 || len(response.Items) != 1 {
		t.Fatalf("expected 1 finished team session, got total=%d len=%d", response.Total, len(response.Items))
	}
	if response.Items[0].RoomID != finishedRoomID {
		t.Fatalf("expected finished room %s, got %s", finishedRoomID, response.Items[0].RoomID)
	}
	if response.Items[0].Role != "VIEWER" {
		t.Fatalf("expected VIEWER role for non-participating owner, got %s", response.Items[0].Role)
	}

	memberReq := httptest.NewRequest(http.MethodGet, "/api/v1/history/teams/"+teamID+"/sessions", nil)
	memberReq.Header.Set("Authorization", "Bearer "+memberToken)
	memberRR := httptest.NewRecorder()
	router.ServeHTTP(memberRR, memberReq)
	if memberRR.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-owner member, got %d: %s", memberRR.Code, memberRR.Body.String())
	}

	outsiderReq := httptest.NewRequest(http.MethodGet, "/api/v1/history/teams/"+teamID+"/sessions", nil)
	outsiderReq.Header.Set("Authorization", "Bearer "+outsiderToken)
	outsiderRR := httptest.NewRecorder()
	router.ServeHTTP(outsiderRR, outsiderReq)
	if outsiderRR.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for outsider, got %d: %s", outsiderRR.Code, outsiderRR.Body.String())
	}
}

func TestListTeamSessions_AllStatusFilterReturnsOK(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createHistoryAccessToken(t, db, "team-all-owner@example.com")

	teamID := seedHistoryTeam(t, db, ownerUserID, "Platform")

	roomID := uuid.NewString()
	createdAt := time.Date(2026, 3, 18, 8, 0, 0, 0, time.UTC)
	lastActivityAt := time.Date(2026, 3, 18, 9, 30, 0, 0, time.UTC)
	seedHistoryRoom(t, db, roomID, "Team All Status Room", ownerUserID, &teamID, "ACTIVE", createdAt, lastActivityAt, nil)
	seedHistoryParticipant(t, db, roomID, ownerUserID, "ADMIN", createdAt)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/teams/"+teamID+"/sessions?status=ALL", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for team status=ALL, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.PaginatedResponse[historydto.SessionListItem]
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Total != 1 || len(response.Items) != 1 {
		t.Fatalf("expected one team session for status=ALL, got total=%d len=%d", response.Total, len(response.Items))
	}
}

func TestGetRoomSummary_AllowsAdminAndAggregatesDetails(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	adminToken, adminUserID := createHistoryAccessToken(t, db, "summary-admin@example.com")
	_, memberUserID := createHistoryAccessToken(t, db, "summary-member@example.com")

	roomID := uuid.NewString()
	roomCreatedAt := time.Date(2026, 3, 15, 9, 0, 0, 0, time.UTC)
	roomActivityAt := time.Date(2026, 3, 15, 11, 0, 0, 0, time.UTC)
	finishedAt := time.Date(2026, 3, 15, 11, 30, 0, 0, time.UTC)
	seedHistoryRoom(t, db, roomID, "Sprint Planning", adminUserID, nil, "FINISHED", roomCreatedAt, roomActivityAt, &finishedAt)

	adminParticipantID := seedHistoryParticipantWithID(t, db, roomID, adminUserID, "ADMIN", roomCreatedAt)
	memberParticipantID := seedHistoryParticipantWithID(t, db, roomID, memberUserID, "MEMBER", roomCreatedAt.Add(2*time.Minute))
	guestParticipantID := seedHistoryGuestParticipantWithID(t, db, roomID, "Guest Estimator", "GUEST", roomCreatedAt.Add(5*time.Minute))
	leftAt := time.Date(2026, 3, 15, 10, 45, 0, 0, time.UTC)
	if _, err := db.ExecContext(context.Background(), `
		UPDATE room_participants
		SET left_at = $2
		WHERE room_participants_id = $1
	`, guestParticipantID, leftAt); err != nil {
		t.Fatalf("failed to set guest left_at: %v", err)
	}

	finalEstimate := "5"
	estimatedTaskID := seedHistoryTaskWithID(
		t,
		db,
		roomID,
		"Backend API",
		"ESTIMATED",
		false,
		&finalEstimate,
		roomCreatedAt.Add(10*time.Minute),
		roomCreatedAt.Add(35*time.Minute),
	)
	votingTaskID := seedHistoryTaskWithID(
		t,
		db,
		roomID,
		"Frontend polish",
		"VOTING",
		true,
		nil,
		roomCreatedAt.Add(40*time.Minute),
		roomCreatedAt.Add(70*time.Minute),
	)

	seedHistoryTaskRound(
		t,
		db,
		estimatedTaskID,
		1,
		"REVEALED",
		[]string{memberParticipantID, guestParticipantID},
		roomCreatedAt.Add(12*time.Minute),
		roomCreatedAt.Add(22*time.Minute),
	)
	seedHistoryTaskRound(
		t,
		db,
		votingTaskID,
		1,
		"ACTIVE",
		[]string{adminParticipantID, memberParticipantID},
		roomCreatedAt.Add(42*time.Minute),
		roomCreatedAt.Add(65*time.Minute),
	)

	seedHistoryVote(t, db, estimatedTaskID, memberParticipantID, 1, "5", roomCreatedAt.Add(15*time.Minute))
	seedHistoryVote(t, db, estimatedTaskID, guestParticipantID, 1, "8", roomCreatedAt.Add(16*time.Minute))
	seedHistoryVote(t, db, votingTaskID, memberParticipantID, 1, "3", roomCreatedAt.Add(50*time.Minute))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/history/rooms/"+roomID+"/summary", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response historydto.RoomSummaryResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode summary response: %v", err)
	}

	if response.Overview.RoomID != roomID {
		t.Fatalf("expected room id %s, got %s", roomID, response.Overview.RoomID)
	}
	if response.Overview.AdminUser.UserID != adminUserID {
		t.Fatalf("expected admin user %s, got %s", adminUserID, response.Overview.AdminUser.UserID)
	}
	if response.Overview.ParticipantsCount != 3 {
		t.Fatalf("expected 3 participants, got %d", response.Overview.ParticipantsCount)
	}
	if response.Overview.TasksCount != 2 {
		t.Fatalf("expected 2 tasks, got %d", response.Overview.TasksCount)
	}
	if response.Overview.EstimatedTasksCount != 1 {
		t.Fatalf("expected 1 estimated task, got %d", response.Overview.EstimatedTasksCount)
	}
	if response.Overview.RoundCount != 2 {
		t.Fatalf("expected 2 rounds, got %d", response.Overview.RoundCount)
	}

	memberFound := false
	guestFound := false
	for _, participant := range response.Participants {
		switch {
		case participant.UserID != nil && *participant.UserID == memberUserID:
			memberFound = true
			if participant.VotesCastCount != 2 {
				t.Fatalf("expected member votesCastCount 2, got %d", participant.VotesCastCount)
			}
			if participant.EstimatedTasksVotedCount != 1 {
				t.Fatalf("expected member estimatedTasksVotedCount 1, got %d", participant.EstimatedTasksVotedCount)
			}
		case participant.GuestName != nil && *participant.GuestName == "Guest Estimator":
			guestFound = true
			if participant.LeftAt == nil {
				t.Fatal("expected guest leftAt to be set")
			}
		}
	}
	if !memberFound {
		t.Fatal("expected member participant in summary")
	}
	if !guestFound {
		t.Fatal("expected guest participant in summary")
	}

	if len(response.Tasks) != 2 {
		t.Fatalf("expected 2 tasks in summary, got %d", len(response.Tasks))
	}

	var estimatedTask historydto.RoomSummaryTask
	var activeTask historydto.RoomSummaryTask
	for _, task := range response.Tasks {
		switch task.TaskID {
		case estimatedTaskID:
			estimatedTask = task
		case votingTaskID:
			activeTask = task
		}
	}

	if estimatedTask.TaskID == "" {
		t.Fatal("expected estimated task in summary")
	}
	if estimatedTask.FinalEstimateValue == nil || *estimatedTask.FinalEstimateValue != finalEstimate {
		t.Fatalf("expected final estimate %s, got %#v", finalEstimate, estimatedTask.FinalEstimateValue)
	}
	if estimatedTask.RoundCount != 1 || len(estimatedTask.Rounds) != 1 {
		t.Fatalf("expected estimated task to have one round, got roundCount=%d len=%d", estimatedTask.RoundCount, len(estimatedTask.Rounds))
	}
	if len(estimatedTask.Rounds[0].Votes) != 2 {
		t.Fatalf("expected 2 revealed votes, got %d", len(estimatedTask.Rounds[0].Votes))
	}

	if activeTask.TaskID == "" {
		t.Fatal("expected active task in summary")
	}
	if activeTask.RoundCount != 1 || len(activeTask.Rounds) != 1 {
		t.Fatalf("expected active task to have one round, got roundCount=%d len=%d", activeTask.RoundCount, len(activeTask.Rounds))
	}
	if activeTask.Rounds[0].Status != "ACTIVE" {
		t.Fatalf("expected active round status ACTIVE, got %s", activeTask.Rounds[0].Status)
	}
	if len(activeTask.Rounds[0].Votes) != 0 {
		t.Fatalf("expected unrevealed round to hide votes, got %d", len(activeTask.Rounds[0].Votes))
	}
}

func TestGetRoomSummary_AllowsTeamOwnerAndRejectsUnauthorizedUsers(t *testing.T) {
	router, db := setupHistoryTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createHistoryAccessToken(t, db, "team-owner@example.com")
	adminToken, adminUserID := createHistoryAccessToken(t, db, "team-room-admin@example.com")
	outsiderToken, _ := createHistoryAccessToken(t, db, "team-outsider@example.com")

	teamID := seedHistoryTeam(t, db, ownerUserID, "Delivery")
	seedHistoryTeamMember(t, db, teamID, adminUserID, "MEMBER")

	roomID := uuid.NewString()
	finishedAt := time.Date(2026, 3, 16, 10, 30, 0, 0, time.UTC)
	seedHistoryRoom(
		t,
		db,
		roomID,
		"Team Room",
		adminUserID,
		&teamID,
		"FINISHED",
		time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 16, 10, 0, 0, 0, time.UTC),
		&finishedAt,
	)
	seedHistoryParticipant(t, db, roomID, adminUserID, "ADMIN", time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC))

	ownerReq := httptest.NewRequest(http.MethodGet, "/api/v1/history/rooms/"+roomID+"/summary", nil)
	ownerReq.Header.Set("Authorization", "Bearer "+ownerToken)
	ownerRR := httptest.NewRecorder()
	router.ServeHTTP(ownerRR, ownerReq)
	if ownerRR.Code != http.StatusOK {
		t.Fatalf("expected 200 for team owner, got %d: %s", ownerRR.Code, ownerRR.Body.String())
	}

	adminReq := httptest.NewRequest(http.MethodGet, "/api/v1/history/rooms/"+roomID+"/summary", nil)
	adminReq.Header.Set("Authorization", "Bearer "+adminToken)
	adminRR := httptest.NewRecorder()
	router.ServeHTTP(adminRR, adminReq)
	if adminRR.Code != http.StatusOK {
		t.Fatalf("expected 200 for room admin, got %d: %s", adminRR.Code, adminRR.Body.String())
	}

	outsiderReq := httptest.NewRequest(http.MethodGet, "/api/v1/history/rooms/"+roomID+"/summary", nil)
	outsiderReq.Header.Set("Authorization", "Bearer "+outsiderToken)
	outsiderRR := httptest.NewRecorder()
	router.ServeHTTP(outsiderRR, outsiderReq)
	if outsiderRR.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for outsider, got %d: %s", outsiderRR.Code, outsiderRR.Body.String())
	}
}
