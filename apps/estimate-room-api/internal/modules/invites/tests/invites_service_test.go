package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupInvitesTest(t *testing.T) (*bun.DB, invites.InvitesService, invitesrepositories.InvitationRepository) {
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

	repo := invitesrepositories.NewInvitationRepository(db)
	svc := invites.NewInvitesService(repo, testutils.TestTokenKey)

	return db, svc, repo
}

func seedInviteTeam(t *testing.T, db *bun.DB, ownerUserID, name string) string {
	t.Helper()

	teamIDValue := name + "-team"
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO teams (team_id, name, owner_user_id)
		VALUES ($1, $2, $3)
	`, teamIDValue, name, ownerUserID)
	if err != nil {
		t.Fatalf("failed to insert team: %v", err)
	}

	_, err = db.ExecContext(context.Background(), `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, 'OWNER')
	`, teamIDValue, ownerUserID)
	if err != nil {
		t.Fatalf("failed to insert team owner: %v", err)
	}

	return teamIDValue
}

func seedInviteRoom(t *testing.T, db *bun.DB, adminUserID, roomID, roomCode string) string {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO rooms (room_id, code, name, admin_user_id, deck)
		VALUES ($1, $2, $3, $4, '{"name":"Fibonacci","kind":"FIBONACCI","values":["1","2","3","5","8"]}'::jsonb)
	`, roomID, roomCode, "Invite Room", adminUserID)
	if err != nil {
		t.Fatalf("failed to insert room: %v", err)
	}

	return roomID
}

func stringPtr(value string) *string {
	return &value
}

func TestCreateInvitation_PersistsAndParsesTeamMemberInvite(t *testing.T) {
	db, svc, repo := setupInvitesTest(t)
	defer db.Close()

	ownerUserID := testutils.SeedUser(t, db, "owner@example.com", "password123")
	invitedUserID := testutils.SeedUser(t, db, "member@example.com", "password123")
	teamID := seedInviteTeam(t, db, ownerUserID, "platform")

	invitation, token, err := svc.CreateInvitation(context.Background(), invites.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindTeamMember,
		TeamID:          &teamID,
		InvitedUserID:   &invitedUserID,
		InvitedEmail:    stringPtr("Member@Example.com"),
		CreatedByUserID: ownerUserID,
	})
	if err != nil {
		t.Fatalf("failed to create invitation: %v", err)
	}

	if invitation.Status != invitesmodels.InvitationStatusActive {
		t.Fatalf("expected ACTIVE status, got %s", invitation.Status)
	}
	if invitation.InvitedEmail == nil || *invitation.InvitedEmail != "member@example.com" {
		t.Fatalf("expected normalized email member@example.com, got %#v", invitation.InvitedEmail)
	}
	if token == "" {
		t.Fatal("expected non-empty invitation token")
	}

	claims, err := svc.ParseInvitationToken(token)
	if err != nil {
		t.Fatalf("failed to parse invitation token: %v", err)
	}

	if claims.InvitationID != invitation.InvitationID {
		t.Fatalf("expected invitation id %s, got %s", invitation.InvitationID, claims.InvitationID)
	}
	if claims.TokenID != invitation.TokenID {
		t.Fatalf("expected token id %s, got %s", invitation.TokenID, claims.TokenID)
	}
	if claims.TeamID == nil || *claims.TeamID != teamID {
		t.Fatalf("expected team id %s, got %#v", teamID, claims.TeamID)
	}

	storedInvitation, err := repo.FindByTokenID(claims.TokenID)
	if err != nil {
		t.Fatalf("failed to load invitation by token id: %v", err)
	}
	if storedInvitation.InvitationID != invitation.InvitationID {
		t.Fatalf("expected stored invitation id %s, got %s", invitation.InvitationID, storedInvitation.InvitationID)
	}

	preview, err := svc.PreviewInvitation(token)
	if err != nil {
		t.Fatalf("failed to preview invitation: %v", err)
	}
	if preview.InvitationID != invitation.InvitationID {
		t.Fatalf("expected preview invitation id %s, got %s", invitation.InvitationID, preview.InvitationID)
	}
}

func TestCreateInvitation_RejectsInvalidPayload(t *testing.T) {
	db, svc, _ := setupInvitesTest(t)
	defer db.Close()

	ownerUserID := testutils.SeedUser(t, db, "owner@example.com", "password123")
	roomID := "room-invalid"
	seedInviteRoom(t, db, ownerUserID, roomID, "room-code-invalid")

	_, _, err := svc.CreateInvitation(context.Background(), invites.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomLink,
		RoomID:          &roomID,
		InvitedEmail:    stringPtr("unexpected@example.com"),
		CreatedByUserID: ownerUserID,
	})
	if !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected bad request, got %v", err)
	}
}

func TestAcceptInvitation_TransitionsToAccepted(t *testing.T) {
	db, svc, repo := setupInvitesTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin@example.com", "password123")
	roomID := seedInviteRoom(t, db, adminUserID, "room-accept", "room-code-accept")

	invitation, token, err := svc.CreateInvitation(context.Background(), invites.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomEmail,
		RoomID:          &roomID,
		InvitedEmail:    stringPtr("invitee@example.com"),
		CreatedByUserID: adminUserID,
	})
	if err != nil {
		t.Fatalf("failed to create room email invite: %v", err)
	}

	acceptedInvitation, err := svc.AcceptInvitation(token)
	if err != nil {
		t.Fatalf("failed to accept invitation: %v", err)
	}
	if acceptedInvitation.Status != invitesmodels.InvitationStatusAccepted {
		t.Fatalf("expected ACCEPTED status, got %s", acceptedInvitation.Status)
	}
	if acceptedInvitation.AcceptedAt == nil {
		t.Fatal("expected acceptedAt to be set")
	}
	if acceptedInvitation.DeclinedAt != nil || acceptedInvitation.RevokedAt != nil {
		t.Fatal("expected only acceptedAt to be set")
	}

	storedInvitation, err := repo.FindByID(invitation.InvitationID)
	if err != nil {
		t.Fatalf("failed to reload invitation: %v", err)
	}
	if storedInvitation.Status != invitesmodels.InvitationStatusAccepted {
		t.Fatalf("expected stored status ACCEPTED, got %s", storedInvitation.Status)
	}

	_, err = svc.AcceptInvitation(token)
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("expected conflict on second accept, got %v", err)
	}
}

func TestDeclineInvitation_TransitionsToDeclined(t *testing.T) {
	db, svc, _ := setupInvitesTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin@example.com", "password123")
	roomID := seedInviteRoom(t, db, adminUserID, "room-decline", "room-code-decline")

	invitation, token, err := svc.CreateInvitation(context.Background(), invites.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomEmail,
		RoomID:          &roomID,
		InvitedEmail:    stringPtr("invitee@example.com"),
		CreatedByUserID: adminUserID,
	})
	if err != nil {
		t.Fatalf("failed to create invitation: %v", err)
	}

	declinedInvitation, err := svc.DeclineInvitation(token)
	if err != nil {
		t.Fatalf("failed to decline invitation: %v", err)
	}
	if declinedInvitation.InvitationID != invitation.InvitationID {
		t.Fatalf("expected invitation id %s, got %s", invitation.InvitationID, declinedInvitation.InvitationID)
	}
	if declinedInvitation.Status != invitesmodels.InvitationStatusDeclined {
		t.Fatalf("expected DECLINED status, got %s", declinedInvitation.Status)
	}
	if declinedInvitation.DeclinedAt == nil {
		t.Fatal("expected declinedAt to be set")
	}
}

func TestRevokeInvitation_TransitionsToRevoked(t *testing.T) {
	db, svc, _ := setupInvitesTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin@example.com", "password123")
	roomID := seedInviteRoom(t, db, adminUserID, "room-revoke", "room-code-revoke")

	invitation, _, err := svc.CreateInvitation(context.Background(), invites.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomLink,
		RoomID:          &roomID,
		CreatedByUserID: adminUserID,
	})
	if err != nil {
		t.Fatalf("failed to create invitation: %v", err)
	}

	revokedInvitation, err := svc.RevokeInvitation(invitation.InvitationID)
	if err != nil {
		t.Fatalf("failed to revoke invitation: %v", err)
	}
	if revokedInvitation.Status != invitesmodels.InvitationStatusRevoked {
		t.Fatalf("expected REVOKED status, got %s", revokedInvitation.Status)
	}
	if revokedInvitation.RevokedAt == nil {
		t.Fatal("expected revokedAt to be set")
	}
}

func TestPreviewInvitation_RejectsTamperedToken(t *testing.T) {
	db, svc, _ := setupInvitesTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin@example.com", "password123")
	roomID := seedInviteRoom(t, db, adminUserID, "room-preview", "room-code-preview")

	_, token, err := svc.CreateInvitation(context.Background(), invites.CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomLink,
		RoomID:          &roomID,
		CreatedByUserID: adminUserID,
	})
	if err != nil {
		t.Fatalf("failed to create invitation: %v", err)
	}

	tamperedToken := token[:len(token)-1] + "x"
	_, err = svc.PreviewInvitation(tamperedToken)
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found for tampered token, got %v", err)
	}
}
