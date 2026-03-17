package invites

import (
	"context"
	"errors"
	"testing"

	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupInternalInvitesTest(t *testing.T) (*bun.DB, *invitesService) {
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
	svc, ok := NewInvitesService(db, repo, testutils.TestTokenKey).(*invitesService)
	if !ok {
		t.Fatal("expected concrete invites service")
	}

	return db, svc
}

func seedInternalInviteRoom(t *testing.T, db *bun.DB, adminUserID, roomID, roomCode string) string {
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

func countActiveRoomParticipants(t *testing.T, db *bun.DB, roomID string) int {
	t.Helper()

	var count int
	err := db.NewSelect().
		TableExpr("room_participants AS rp").
		ColumnExpr("COUNT(*)").
		Where("rp.room_id = ?", roomID).
		Where("rp.left_at IS NULL").
		Scan(context.Background(), &count)
	if err != nil {
		t.Fatalf("failed to count room participants: %v", err)
	}

	return count
}

func TestAcceptRoomLinkInvitation_RejectsRevokedStaleInvitationForRegisteredUser(t *testing.T) {
	db, svc := setupInternalInvitesTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "room-link-admin@example.com", "password123")
	memberUserID := testutils.SeedUser(t, db, "room-link-member@example.com", "password123")
	roomID := seedInternalInviteRoom(t, db, adminUserID, "room-link-stale-member", "room-link-stale-member")

	invitation, token, err := svc.CreateInvitation(context.Background(), CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomLink,
		RoomID:          &roomID,
		CreatedByUserID: adminUserID,
	})
	if err != nil {
		t.Fatalf("failed to create room link invitation: %v", err)
	}

	staleInvitation, err := svc.PreviewInvitation(context.Background(), token)
	if err != nil {
		t.Fatalf("failed to preview invitation: %v", err)
	}

	if _, err := svc.RevokeInvitation(context.Background(), invitation.InvitationID, adminUserID); err != nil {
		t.Fatalf("failed to revoke invitation: %v", err)
	}

	result, err := svc.acceptRoomLinkInvitation(context.Background(), staleInvitation, memberUserID, nil)
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("expected conflict for revoked room link invitation, got result=%#v err=%v", result, err)
	}

	if count := countActiveRoomParticipants(t, db, roomID); count != 0 {
		t.Fatalf("expected no active room participants after rejected join, got %d", count)
	}
}

func TestAcceptRoomLinkInvitation_RejectsRevokedStaleInvitationForGuest(t *testing.T) {
	db, svc := setupInternalInvitesTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "room-link-guest-admin@example.com", "password123")
	roomID := seedInternalInviteRoom(t, db, adminUserID, "room-link-stale-guest", "room-link-stale-guest")

	invitation, token, err := svc.CreateInvitation(context.Background(), CreateInvitationInput{
		Kind:            invitesmodels.InvitationKindRoomLink,
		RoomID:          &roomID,
		CreatedByUserID: adminUserID,
	})
	if err != nil {
		t.Fatalf("failed to create room link invitation: %v", err)
	}

	staleInvitation, err := svc.PreviewInvitation(context.Background(), token)
	if err != nil {
		t.Fatalf("failed to preview invitation: %v", err)
	}

	if _, err := svc.RevokeInvitation(context.Background(), invitation.InvitationID, adminUserID); err != nil {
		t.Fatalf("failed to revoke invitation: %v", err)
	}

	guestName := "Guest Joiner"
	result, err := svc.acceptRoomLinkInvitation(context.Background(), staleInvitation, "", &guestName)
	if !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("expected conflict for revoked guest room link invitation, got result=%#v err=%v", result, err)
	}

	if count := countActiveRoomParticipants(t, db, roomID); count != 0 {
		t.Fatalf("expected no active room participants after rejected guest join, got %d", count)
	}
}
