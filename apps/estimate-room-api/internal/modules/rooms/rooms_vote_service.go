package rooms

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsVoteService interface {
	SetCurrentTask(roomID, taskID, userID string, eligibleParticipantIDs []string) (*roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskRoundModel, error)
	CastVote(roomID string, participant *roomsmodels.RoomParticipantModel, value string) (*CastVoteResult, error)
	RevealCurrentRound(roomID, userID string) (*RevealVotesResult, error)
	StartNextRound(roomID, userID string, eligibleParticipantIDs []string) (*roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskRoundModel, error)
	FinalizeCurrentTask(roomID, userID, value string) (*roomsmodels.RoomTaskModel, error)
	FinalizeTask(roomID, taskID, userID, value string) (*roomsmodels.RoomTaskModel, error)
}

type CastVoteResult struct {
	Task                   *roomsmodels.RoomTaskModel
	Round                  *roomsmodels.RoomTaskRoundModel
	VotedParticipantIDs    []string
	EligibleParticipantIDs []string
	AllVotesCast           bool
}

type RevealVotesResult struct {
	Task     *roomsmodels.RoomTaskModel
	Round    *roomsmodels.RoomTaskRoundModel
	Votes    []*roomsmodels.RoomVoteModel
	AllVoted bool
}

type roomsVoteService struct {
	roomsRepo       roomsrepositories.RoomsRepository
	taskRepo        roomsrepositories.RoomTaskRepository
	voteRepo        roomsrepositories.RoomVoteRepository
	roundRepo       roomsrepositories.RoomTaskRoundRepository
	participantRepo roomsrepositories.RoomParticipantRepository
	logger          *slog.Logger
}

func NewRoomsVoteService(
	roomsRepo roomsrepositories.RoomsRepository,
	taskRepo roomsrepositories.RoomTaskRepository,
	voteRepo roomsrepositories.RoomVoteRepository,
	roundRepo roomsrepositories.RoomTaskRoundRepository,
	participantRepo roomsrepositories.RoomParticipantRepository,
) RoomsVoteService {
	return &roomsVoteService{
		roomsRepo:       roomsRepo,
		taskRepo:        taskRepo,
		voteRepo:        voteRepo,
		roundRepo:       roundRepo,
		participantRepo: participantRepo,
		logger:          logger.L().With(slog.String("service", "rooms-votes")),
	}
}

func (s *roomsVoteService) SetCurrentTask(roomID, taskID, userID string, eligibleParticipantIDs []string) (*roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskRoundModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, nil, nil, err
	}

	task, previousTask, err := s.taskRepo.SetCurrentVotingTask(roomID, taskID)
	if err != nil {
		return nil, nil, nil, err
	}

	round, err := s.roundRepo.GetOrCreateCurrent(task.TaskID, normalizeParticipantIDs(eligibleParticipantIDs))
	if err != nil {
		return nil, nil, nil, err
	}

	return task, previousTask, round, nil
}

func (s *roomsVoteService) CastVote(roomID string, participant *roomsmodels.RoomParticipantModel, value string) (*CastVoteResult, error) {
	if participant == nil {
		return nil, apperrors.ErrUnauthorized
	}
	if !isVotingParticipantRole(participant.Role) {
		return nil, apperrors.ErrForbidden
	}

	task, err := s.taskRepo.FindCurrentVotingTask(roomID)
	if err != nil {
		return nil, err
	}

	round, err := s.roundRepo.GetOrCreateCurrent(task.TaskID, nil)
	if err != nil {
		return nil, err
	}
	if round.Status != roomsmodels.RoomTaskRoundStatusActive {
		return nil, fmt.Errorf("%w: round already revealed", apperrors.ErrBadRequest)
	}
	if !containsParticipantID(round.EligibleParticipantIDs, participant.RoomParticipantID) {
		return nil, apperrors.ErrForbidden
	}

	room, err := s.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	trimmedValue := strings.TrimSpace(value)
	if !isDeckValueAllowed(room.Deck.Values, trimmedValue) {
		return nil, fmt.Errorf("%w: vote value is not in the room deck", apperrors.ErrBadRequest)
	}

	if _, err := s.voteRepo.Upsert(task.TaskID, participant.RoomParticipantID, round.RoundNumber, trimmedValue); err != nil {
		return nil, err
	}

	votes, err := s.voteRepo.ListByTaskAndRound(task.TaskID, round.RoundNumber)
	if err != nil {
		return nil, err
	}

	votedParticipantIDs := filterParticipantIDs(uniqueSortedParticipantIDs(votes), round.EligibleParticipantIDs)
	allVotesCast := len(round.EligibleParticipantIDs) > 0 && sameParticipantIDs(votedParticipantIDs, round.EligibleParticipantIDs)

	return &CastVoteResult{
		Task:                   task,
		Round:                  round,
		VotedParticipantIDs:    votedParticipantIDs,
		EligibleParticipantIDs: append([]string(nil), round.EligibleParticipantIDs...),
		AllVotesCast:           allVotesCast,
	}, nil
}

func (s *roomsVoteService) RevealCurrentRound(roomID, userID string) (*RevealVotesResult, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.FindCurrentVotingTask(roomID)
	if err != nil {
		return nil, err
	}

	round, err := s.roundRepo.GetOrCreateCurrent(task.TaskID, nil)
	if err != nil {
		return nil, err
	}

	votes, err := s.voteRepo.ListByTaskAndRound(task.TaskID, round.RoundNumber)
	if err != nil {
		return nil, err
	}

	allVoted := len(round.EligibleParticipantIDs) > 0 &&
		sameParticipantIDs(filterParticipantIDs(uniqueSortedParticipantIDs(votes), round.EligibleParticipantIDs), round.EligibleParticipantIDs)

	if round.Status == roomsmodels.RoomTaskRoundStatusActive {
		round, err = s.roundRepo.MarkRevealed(task.TaskID, round.RoundNumber)
		if err != nil {
			return nil, err
		}
	}

	return &RevealVotesResult{
		Task:     task,
		Round:    round,
		Votes:    votes,
		AllVoted: allVoted,
	}, nil
}

func (s *roomsVoteService) StartNextRound(roomID, userID string, eligibleParticipantIDs []string) (*roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskRoundModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, nil, err
	}

	task, err := s.taskRepo.FindCurrentVotingTask(roomID)
	if err != nil {
		return nil, nil, err
	}

	currentRound, err := s.roundRepo.GetOrCreateCurrent(task.TaskID, nil)
	if err != nil {
		return nil, nil, err
	}
	if currentRound.Status != roomsmodels.RoomTaskRoundStatusRevealed {
		return nil, nil, fmt.Errorf("%w: current round must be revealed before starting a new round", apperrors.ErrBadRequest)
	}

	round, err := s.roundRepo.Advance(task.TaskID, normalizeParticipantIDs(eligibleParticipantIDs))
	if err != nil {
		return nil, nil, err
	}

	return task, round, nil
}

func (s *roomsVoteService) FinalizeCurrentTask(roomID, userID, value string) (*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.FindCurrentVotingTask(roomID)
	if err != nil {
		return nil, err
	}

	return s.finalizeTask(roomID, task, value)
}

func (s *roomsVoteService) FinalizeTask(roomID, taskID, userID, value string) (*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.FindByID(roomID, taskID)
	if err != nil {
		return nil, err
	}

	return s.finalizeTask(roomID, task, value)
}

func (s *roomsVoteService) finalizeTask(roomID string, task *roomsmodels.RoomTaskModel, value string) (*roomsmodels.RoomTaskModel, error) {
	room, err := s.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return nil, fmt.Errorf("%w: final estimate value is required", apperrors.ErrBadRequest)
	}
	if task.Status == "SKIPPED" {
		return nil, fmt.Errorf("%w: skipped tasks cannot be finalized", apperrors.ErrBadRequest)
	}
	if !isDeckValueAllowed(room.Deck.Values, trimmedValue) {
		return nil, fmt.Errorf("%w: final estimate must be from the room deck", apperrors.ErrBadRequest)
	}
	if task.IsActive {
		round, err := s.roundRepo.GetOrCreateCurrent(task.TaskID, nil)
		if err != nil {
			return nil, err
		}
		if round.Status != roomsmodels.RoomTaskRoundStatusRevealed {
			return nil, fmt.Errorf("%w: current round must be revealed before finalizing", apperrors.ErrBadRequest)
		}
	}

	task.Status = "ESTIMATED"
	task.IsActive = false
	task.FinalEstimateValue = &trimmedValue
	return s.taskRepo.Update(roomID, task)
}

func (s *roomsVoteService) ensureRoomAdmin(roomID, userID string) (*roomsmodels.RoomsModel, error) {
	room, err := s.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	participant, err := s.participantRepo.FindActiveByUserID(roomID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}
		return nil, err
	}
	if participant.Role != roomsmodels.RoomParticipantRoleAdmin {
		return nil, apperrors.ErrForbidden
	}
	if room.AdminUserID != userID {
		return nil, apperrors.ErrForbidden
	}

	return room, nil
}

func isVotingParticipantRole(role roomsmodels.RoomParticipantRole) bool {
	return role == roomsmodels.RoomParticipantRoleMember || role == roomsmodels.RoomParticipantRoleGuest
}

func normalizeParticipantIDs(ids []string) []string {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			set[trimmed] = struct{}{}
		}
	}

	normalized := make([]string, 0, len(set))
	for id := range set {
		normalized = append(normalized, id)
	}
	sort.Strings(normalized)
	return normalized
}

func containsParticipantID(ids []string, target string) bool {
	trimmedTarget := strings.TrimSpace(target)
	if trimmedTarget == "" {
		return false
	}

	for _, id := range ids {
		if strings.TrimSpace(id) == trimmedTarget {
			return true
		}
	}

	return false
}

func filterParticipantIDs(ids []string, allowed []string) []string {
	filtered := make([]string, 0, len(ids))
	for _, id := range ids {
		if containsParticipantID(allowed, id) {
			filtered = append(filtered, id)
		}
	}
	return normalizeParticipantIDs(filtered)
}

func sameParticipantIDs(left []string, right []string) bool {
	normalizedLeft := normalizeParticipantIDs(left)
	normalizedRight := normalizeParticipantIDs(right)
	if len(normalizedLeft) != len(normalizedRight) {
		return false
	}
	for i := range normalizedLeft {
		if normalizedLeft[i] != normalizedRight[i] {
			return false
		}
	}
	return true
}
