package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

type failingInvitesService struct {
	err error
}

func (s *failingInvitesService) CreateInvitation(
	ctx context.Context,
	input invites.CreateInvitationInput,
) (*invitesmodels.InvitationModel, string, error) {
	return nil, "", s.err
}

func (s *failingInvitesService) CreateInvitationWithDB(
	ctx context.Context,
	db bun.IDB,
	input invites.CreateInvitationInput,
) (*invitesmodels.InvitationModel, string, error) {
	return nil, "", s.err
}

func (s *failingInvitesService) ParseInvitationToken(token string) (*invites.InvitationTokenClaims, error) {
	return nil, s.err
}

func (s *failingInvitesService) PreviewInvitation(
	ctx context.Context,
	token string,
) (*invitesmodels.InvitationModel, error) {
	return nil, s.err
}

func (s *failingInvitesService) AcceptInvitation(
	ctx context.Context,
	token, actorUserID string,
	guestName *string,
) (*invites.AcceptInvitationResult, error) {
	return nil, s.err
}

func (s *failingInvitesService) DeclineInvitation(
	ctx context.Context,
	token, actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	return nil, s.err
}

func (s *failingInvitesService) RevokeInvitation(
	ctx context.Context,
	invitationID, actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	return nil, s.err
}

func (s *failingInvitesService) ValidateGuestRoomAccess(
	roomID, guestToken string,
) (*roomsmodels.RoomParticipantModel, error) {
	return nil, s.err
}

func TestCreateRoom_RollsBackRoomAndParticipantWhenInviteCreationFails(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

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

	adminUserID := testutils.SeedUser(t, db, "admin@example.com", "password123")
	service := rooms.NewRoomsService(
		db,
		roomsrepositories.NewRoomsRepository(db),
		roomsrepositories.NewRoomParticipantRepository(db),
		teamsrepositories.NewTeamRepository(db),
		teamsrepositories.NewTeamMemberRepository(db),
		usersrepositories.NewUserRepository(db),
		&failingInvitesService{err: errors.New("invite failure")},
		nil,
	)

	_, err = service.CreateRoom(context.Background(), rooms.CreateRoomInput{
		Name:         "Atomic Room",
		AdminUserID:  adminUserID,
		InviteEmails: []string{"invitee@example.com"},
	})
	if err == nil {
		t.Fatal("expected create room to fail when invite creation fails")
	}

	var roomCount int
	if err := db.NewSelect().Table("rooms").ColumnExpr("COUNT(*)").Scan(context.Background(), &roomCount); err != nil {
		t.Fatalf("failed to count rooms: %v", err)
	}
	if roomCount != 0 {
		t.Fatalf("expected no rooms after rollback, got %d", roomCount)
	}

	var participantCount int
	if err := db.NewSelect().
		Table("room_participants").
		ColumnExpr("COUNT(*)").
		Scan(context.Background(), &participantCount); err != nil {
		t.Fatalf("failed to count room participants: %v", err)
	}
	if participantCount != 0 {
		t.Fatalf("expected no room participants after rollback, got %d", participantCount)
	}

	var invitationCount int
	if err := db.NewSelect().
		Table("invitations").
		ColumnExpr("COUNT(*)").
		Scan(context.Background(), &invitationCount); err != nil {
		t.Fatalf("failed to count invitations: %v", err)
	}
	if invitationCount != 0 {
		t.Fatalf("expected no invitations after rollback, got %d", invitationCount)
	}
}
