package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupRoomsVoteServiceTest(t *testing.T) (*bun.DB, rooms.RoomsVoteService, roomsrepositories.RoomParticipantRepository) {
	t.Helper()

	db := testutils.SetupTestDB(t)

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
			invitations,
			votes,
			task_rounds,
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

	roomsRepo := roomsrepositories.NewRoomsRepository(db)
	participantRepo := roomsrepositories.NewRoomParticipantRepository(db)
	expiryService := rooms.NewRoomsExpiryService(db, roomsRepo, nil, nil)
	voteService := rooms.NewRoomsVoteService(
		roomsRepo,
		roomsrepositories.NewRoomTaskRepository(db),
		roomsrepositories.NewRoomVoteRepository(db),
		roomsrepositories.NewRoomTaskRoundRepository(db),
		participantRepo,
		expiryService,
	)

	return db, voteService, participantRepo
}

func TestRoomsVoteService_CastVoteRejectsValueOutsideDeck(t *testing.T) {
	db, voteService, participantRepo := setupRoomsVoteServiceTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin-vote@example.com", "password123")
	memberUserID := testutils.SeedUser(t, db, "member-vote@example.com", "password123")

	roomID := seedRoom(t, db, adminUserID)
	memberParticipantID := seedMemberParticipant(t, db, roomID, memberUserID)
	taskID := seedTask(t, db, roomID, "Vote bounds")

	if _, _, _, err := voteService.SetCurrentTask(roomID, taskID, adminUserID, []string{memberParticipantID}); err != nil {
		t.Fatalf("failed to set current task: %v", err)
	}

	participant, err := participantRepo.FindActiveByUserID(roomID, memberUserID)
	if err != nil {
		t.Fatalf("failed to load member participant: %v", err)
	}

	_, err = voteService.CastVote(roomID, participant, "13")
	if !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest for invalid deck vote, got %v", err)
	}
}

func TestRoomsVoteService_StartNextRoundRequiresRevealedRound(t *testing.T) {
	db, voteService, _ := setupRoomsVoteServiceTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin-round@example.com", "password123")
	memberUserID := testutils.SeedUser(t, db, "member-round@example.com", "password123")

	roomID := seedRoom(t, db, adminUserID)
	memberParticipantID := seedMemberParticipant(t, db, roomID, memberUserID)
	taskID := seedTask(t, db, roomID, "Round guard")

	if _, _, _, err := voteService.SetCurrentTask(roomID, taskID, adminUserID, []string{memberParticipantID}); err != nil {
		t.Fatalf("failed to set current task: %v", err)
	}

	_, _, err := voteService.StartNextRound(roomID, adminUserID, []string{memberParticipantID})
	if !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest when advancing unrevealed round, got %v", err)
	}
}

func TestRoomsVoteService_FinalizeCurrentTaskRequiresRevealedRound(t *testing.T) {
	db, voteService, _ := setupRoomsVoteServiceTest(t)
	defer db.Close()

	adminUserID := testutils.SeedUser(t, db, "admin-finalize@example.com", "password123")
	memberUserID := testutils.SeedUser(t, db, "member-finalize@example.com", "password123")

	roomID := seedRoom(t, db, adminUserID)
	memberParticipantID := seedMemberParticipant(t, db, roomID, memberUserID)
	taskID := seedTask(t, db, roomID, "Finalize guard")

	if _, _, _, err := voteService.SetCurrentTask(roomID, taskID, adminUserID, []string{memberParticipantID}); err != nil {
		t.Fatalf("failed to set current task: %v", err)
	}

	_, err := voteService.FinalizeCurrentTask(roomID, adminUserID, "5")
	if !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest when finalizing before reveal, got %v", err)
	}
}
