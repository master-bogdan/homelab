package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	invitesmodule "github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupInvitesControllerTest(t *testing.T) (*chi.Mux, *bun.DB) {
	t.Helper()

	db := testutils.SetupTestDB(t)

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
			invitations,
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
	authService := oauth2.NewOauth2SessionAuthServiceFromDB(testutils.TestTokenKey, db)

	router.Route("/api/v1", func(r chi.Router) {
		invitesmodule.NewInvitesModule(invitesmodule.InvitesModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
		})
	})

	return router, db
}

func createInviteAccessToken(t *testing.T, db *bun.DB, email string) (string, string) {
	t.Helper()

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, email, "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-invites")

	svc := testutils.NewOauth2Service(db)
	tokens, err := svc.GenerateTokenPair(context.Background(), userID, clientID, sessionID, []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	return tokens.AccessToken, userID
}

func seedInviteControllerTeam(t *testing.T, db *bun.DB, ownerUserID, teamID, teamName string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO teams (team_id, name, owner_user_id)
		VALUES ($1, $2, $3)
	`, teamID, teamName, ownerUserID)
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
}

func createTeamInviteToken(
	t *testing.T,
	db *bun.DB,
	teamID, ownerUserID, invitedUserID, invitedEmail string,
) string {
	t.Helper()

	repo := invitesrepositories.NewInvitationRepository(db)
	svc := invitesmodule.NewInvitesService(db, repo, testutils.TestTokenKey)
	_, token, err := svc.CreateInvitation(context.Background(), invitesmodule.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindTeamMember,
		TeamID:          &teamID,
		InvitedUserID:   &invitedUserID,
		InvitedEmail:    &invitedEmail,
		CreatedByUserID: ownerUserID,
	})
	if err != nil {
		t.Fatalf("failed to create invitation: %v", err)
	}

	return token
}

func TestPreviewInvitation_ReturnsTeamInviteWithoutAuth(t *testing.T) {
	router, db := setupInvitesControllerTest(t)
	defer db.Close()

	_, ownerUserID := createInviteAccessToken(t, db, "owner@example.com")
	_, invitedUserID := createInviteAccessToken(t, db, "member@example.com")
	teamID := "team-preview"
	seedInviteControllerTeam(t, db, ownerUserID, teamID, "Preview Team")
	token := createTeamInviteToken(t, db, teamID, ownerUserID, invitedUserID, "member@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/invites/"+token, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response invitesdto.InvitationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Kind != "TEAM_MEMBER" {
		t.Fatalf("expected TEAM_MEMBER kind, got %s", response.Kind)
	}
	if response.TeamID == nil || *response.TeamID != teamID {
		t.Fatalf("expected team id %s, got %#v", teamID, response.TeamID)
	}
}

func TestAcceptInvitation_AddsTeamMemberForInvitedUser(t *testing.T) {
	router, db := setupInvitesControllerTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createInviteAccessToken(t, db, "owner@example.com")
	_ = ownerToken
	memberToken, invitedUserID := createInviteAccessToken(t, db, "member@example.com")
	teamID := "team-accept"
	seedInviteControllerTeam(t, db, ownerUserID, teamID, "Accept Team")
	token := createTeamInviteToken(t, db, teamID, ownerUserID, invitedUserID, "member@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/invites/"+token+"/accept", nil)
	req.Header.Set("Authorization", "Bearer "+memberToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response invitesdto.InvitationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Status != "ACCEPTED" {
		t.Fatalf("expected ACCEPTED status, got %s", response.Status)
	}

	var memberCount int
	err := db.NewSelect().
		Table("team_members").
		ColumnExpr("COUNT(*)").
		Where("team_id = ?", teamID).
		Where("user_id = ?", invitedUserID).
		Scan(context.Background(), &memberCount)
	if err != nil {
		t.Fatalf("failed to count team members: %v", err)
	}
	if memberCount != 1 {
		t.Fatalf("expected invited user to be added to team, found %d rows", memberCount)
	}
}

func TestAcceptInvitation_RejectsDifferentAuthenticatedUser(t *testing.T) {
	router, db := setupInvitesControllerTest(t)
	defer db.Close()

	_, ownerUserID := createInviteAccessToken(t, db, "owner@example.com")
	_, invitedUserID := createInviteAccessToken(t, db, "member@example.com")
	otherToken, _ := createInviteAccessToken(t, db, "other@example.com")
	teamID := "team-forbidden"
	seedInviteControllerTeam(t, db, ownerUserID, teamID, "Forbidden Team")
	token := createTeamInviteToken(t, db, teamID, ownerUserID, invitedUserID, "member@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/invites/"+token+"/accept", nil)
	req.Header.Set("Authorization", "Bearer "+otherToken)

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

func TestDeclineInvitation_TransitionsToDeclinedForInvitedUser(t *testing.T) {
	router, db := setupInvitesControllerTest(t)
	defer db.Close()

	_, ownerUserID := createInviteAccessToken(t, db, "owner@example.com")
	memberToken, invitedUserID := createInviteAccessToken(t, db, "member@example.com")
	teamID := "team-decline"
	seedInviteControllerTeam(t, db, ownerUserID, teamID, "Decline Team")
	token := createTeamInviteToken(t, db, teamID, ownerUserID, invitedUserID, "member@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/invites/"+token+"/decline", nil)
	req.Header.Set("Authorization", "Bearer "+memberToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response invitesdto.InvitationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Status != "DECLINED" {
		t.Fatalf("expected DECLINED status, got %s", response.Status)
	}
}

func TestRevokeInvitation_AllowsTeamOwner(t *testing.T) {
	router, db := setupInvitesControllerTest(t)
	defer db.Close()

	ownerToken, ownerUserID := createInviteAccessToken(t, db, "owner@example.com")
	_, invitedUserID := createInviteAccessToken(t, db, "member@example.com")
	teamID := "team-revoke"
	seedInviteControllerTeam(t, db, ownerUserID, teamID, "Revoke Team")
	token := createTeamInviteToken(t, db, teamID, ownerUserID, invitedUserID, "member@example.com")

	previewReq := httptest.NewRequest(http.MethodGet, "/api/v1/invites/"+token, nil)
	previewRR := httptest.NewRecorder()
	router.ServeHTTP(previewRR, previewReq)

	var preview invitesdto.InvitationResponse
	if err := json.NewDecoder(previewRR.Body).Decode(&preview); err != nil {
		t.Fatalf("failed to decode preview response: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/invites/"+preview.InvitationID+"/revoke", nil)
	req.Header.Set("Authorization", "Bearer "+ownerToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var response invitesdto.InvitationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Status != "REVOKED" {
		t.Fatalf("expected REVOKED status, got %s", response.Status)
	}
}
