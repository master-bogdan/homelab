package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	invitesmodule "github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	teams "github.com/master-bogdan/estimate-room-api/internal/modules/teams"
	teamsdto "github.com/master-bogdan/estimate-room-api/internal/modules/teams/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupTeamsTest(t *testing.T) (*chi.Mux, *bun.DB) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	router := chi.NewRouter()
	authService := oauth2.NewAuthServiceFromDB(testutils.TestTokenKey, db)
	userRepo := usersrepositories.NewUserRepository(db)
	userService := users.NewUsersService(userRepo)

	router.Route("/api/v1", func(r chi.Router) {
		invitesModule := invitesmodule.NewInvitesModule(invitesmodule.InvitesModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
		})

		teams.NewTeamsModule(teams.TeamsModuleDeps{
			Router:         r,
			DB:             db,
			AuthService:    authService,
			UserService:    userService,
			InvitesService: invitesModule.Service,
		})
	})

	return router, db
}

func createTeamsAccessToken(t *testing.T, db *bun.DB, email string) (string, string) {
	t.Helper()

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, email, "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-teams")

	svc := testutils.NewOauth2Service(db)
	tokens, err := svc.GenerateTokenPair(context.Background(), userID, clientID, sessionID, []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	return tokens.AccessToken, userID
}

func seedTeam(t *testing.T, db *bun.DB, ownerUserID, name string) string {
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
		t.Fatalf("failed to insert owner membership: %v", err)
	}

	return teamID
}

func seedTeamMember(t *testing.T, db *bun.DB, teamID, userID, role string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
	`, teamID, userID, role)
	if err != nil {
		t.Fatalf("failed to insert team member: %v", err)
	}
}

func TestCreateTeam_CreatesOwnerMembership(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	accessToken, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams/", bytes.NewReader([]byte(`{"name":"Platform Team"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response teamsdto.TeamDetailResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.TeamID == "" {
		t.Fatal("expected team id")
	}
	if response.Name != "Platform Team" {
		t.Fatalf("expected team name, got %s", response.Name)
	}
	if response.OwnerUserID != ownerUserID {
		t.Fatalf("expected owner user id %s, got %s", ownerUserID, response.OwnerUserID)
	}
	if len(response.Members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(response.Members))
	}
	if response.Members[0].UserID != ownerUserID {
		t.Fatalf("expected owner member user id %s, got %s", ownerUserID, response.Members[0].UserID)
	}
	if response.Members[0].Role != "OWNER" {
		t.Fatalf("expected OWNER role, got %s", response.Members[0].Role)
	}
}

func TestListTeams_ReturnsOnlyCurrentUsersTeams(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	accessToken, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	_, outsiderUserID := createTeamsAccessToken(t, db, "outsider@example.com")

	firstTeamID := seedTeam(t, db, ownerUserID, "Alpha")
	secondTeamID := seedTeam(t, db, ownerUserID, "Beta")
	seedTeam(t, db, outsiderUserID, "Gamma")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teams/", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response []teamsdto.TeamSummaryResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(response))
	}

	teamIDs := map[string]bool{
		response[0].TeamID: true,
		response[1].TeamID: true,
	}
	if !teamIDs[firstTeamID] || !teamIDs[secondTeamID] {
		t.Fatalf("expected list to contain teams %s and %s, got %#v", firstTeamID, secondTeamID, response)
	}
}

func TestGetTeam_ReturnsMembersForTeamMember(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	accessToken, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	_, memberUserID := createTeamsAccessToken(t, db, "member@example.com")

	teamID := seedTeam(t, db, ownerUserID, "Product")
	seedTeamMember(t, db, teamID, memberUserID, "MEMBER")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teams/"+teamID, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response teamsdto.TeamDetailResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.TeamID != teamID {
		t.Fatalf("expected team id %s, got %s", teamID, response.TeamID)
	}
	if len(response.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(response.Members))
	}

	memberIDs := map[string]bool{}
	for _, member := range response.Members {
		memberIDs[member.UserID] = true
	}
	if !memberIDs[ownerUserID] || !memberIDs[memberUserID] {
		t.Fatalf("expected members %s and %s, got %#v", ownerUserID, memberUserID, response.Members)
	}
}

func TestRemoveMember_RemovesNonOwnerMember(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	accessToken, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	_, memberUserID := createTeamsAccessToken(t, db, "member@example.com")

	teamID := seedTeam(t, db, ownerUserID, "Backend")
	seedTeamMember(t, db, teamID, memberUserID, "MEMBER")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/teams/"+teamID+"/members/"+memberUserID, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var memberCount int
	err := db.NewSelect().
		Table("team_members").
		ColumnExpr("COUNT(*)").
		Where("team_id = ?", teamID).
		Where("user_id = ?", memberUserID).
		Scan(context.Background(), &memberCount)
	if err != nil {
		t.Fatalf("failed to count team members: %v", err)
	}
	if memberCount != 0 {
		t.Fatalf("expected member to be removed, found %d rows", memberCount)
	}
}

func TestRemoveMember_NonOwnerIsForbidden(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	_, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	memberToken, memberUserID := createTeamsAccessToken(t, db, "member@example.com")
	_, anotherUserID := createTeamsAccessToken(t, db, "another@example.com")

	teamID := seedTeam(t, db, ownerUserID, "Infra")
	seedTeamMember(t, db, teamID, memberUserID, "MEMBER")
	seedTeamMember(t, db, teamID, anotherUserID, "MEMBER")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/teams/"+teamID+"/members/"+anotherUserID, nil)
	req.Header.Set("Authorization", "Bearer "+memberToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d: %s", rr.Code, rr.Body.String())
	}

	var httpErr apperrors.HttpError
	if err := json.NewDecoder(rr.Body).Decode(&httpErr); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if httpErr.Status != http.StatusForbidden {
		t.Fatalf("expected error status 403, got %d", httpErr.Status)
	}
}

func TestCreateInvites_OwnerCanInviteRegisteredUsers(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	_, _ = createTeamsAccessToken(t, db, "member@example.com")
	teamID := seedTeam(t, db, ownerUserID, "Invites Team")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams/"+teamID+"/invites", bytes.NewReader([]byte(`{"emails":["member@example.com"]}`)))
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response []invitesdto.InvitationWithTokenResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("expected 1 invitation, got %d", len(response))
	}
	if response[0].Kind != "TEAM_MEMBER" {
		t.Fatalf("expected TEAM_MEMBER kind, got %s", response[0].Kind)
	}
	if response[0].Status != "ACTIVE" {
		t.Fatalf("expected ACTIVE status, got %s", response[0].Status)
	}
	if response[0].Token == "" {
		t.Fatal("expected invitation token")
	}
}

func TestCreateInvites_RejectsUnknownEmail(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	teamID := seedTeam(t, db, ownerUserID, "Unknown Invite Team")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams/"+teamID+"/invites", bytes.NewReader([]byte(`{"emails":["unknown@example.com"]}`)))
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCreateInvites_NonOwnerIsForbidden(t *testing.T) {
	router, db := setupTeamsTest(t)
	defer db.Close()

	_, ownerUserID := createTeamsAccessToken(t, db, "owner@example.com")
	memberToken, memberUserID := createTeamsAccessToken(t, db, "member@example.com")
	_, _ = createTeamsAccessToken(t, db, "invitee@example.com")
	teamID := seedTeam(t, db, ownerUserID, "Owner Only Team")
	seedTeamMember(t, db, teamID, memberUserID, "MEMBER")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams/"+teamID+"/invites", bytes.NewReader([]byte(`{"emails":["invitee@example.com"]}`)))
	req.Header.Set("Authorization", "Bearer "+memberToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d: %s", rr.Code, rr.Body.String())
	}
}
