package rooms

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsTaskService interface {
	CreateTask(roomID, userID string, input CreateTaskInput) (*roomsmodels.RoomTaskModel, error)
	ListTasks(roomID, userID string) ([]*roomsmodels.RoomTaskModel, error)
	GetTask(roomID, taskID, userID string) (*roomsmodels.RoomTaskModel, error)
	UpdateTask(roomID, taskID, userID string, input UpdateTaskInput) (*roomsmodels.RoomTaskModel, error)
	DeleteTask(roomID, taskID, userID string) error
}

type roomsTaskService struct {
	roomsRepo       roomsrepositories.RoomsRepository
	taskRepo        roomsrepositories.RoomTaskRepository
	voteService     RoomsVoteService
	participantRepo roomsrepositories.RoomParticipantRepository
	expiryService   RoomsExpiryService
	logger          *slog.Logger
}

type CreateTaskInput struct {
	Title       string
	Description string
	ExternalKey string
}

type UpdateTaskInput struct {
	Title              *string
	Description        *string
	ExternalKey        *string
	Status             *string
	IsActive           *bool
	FinalEstimateValue *string
}

func NewRoomsTaskService(
	roomsRepo roomsrepositories.RoomsRepository,
	taskRepo roomsrepositories.RoomTaskRepository,
	voteService RoomsVoteService,
	participantRepo roomsrepositories.RoomParticipantRepository,
	expiryService RoomsExpiryService,
) RoomsTaskService {
	return &roomsTaskService{
		roomsRepo:       roomsRepo,
		taskRepo:        taskRepo,
		voteService:     voteService,
		participantRepo: participantRepo,
		expiryService:   expiryService,
		logger:          logger.L().With(slog.String("service", "rooms-tasks")),
	}
}

func (s *roomsTaskService) CreateTask(roomID, userID string, input CreateTaskInput) (*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureActiveRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, apperrors.ErrBadRequest
	}

	var description *string
	if input.Description != "" {
		description = &input.Description
	}

	var externalKey *string
	if input.ExternalKey != "" {
		externalKey = &input.ExternalKey
	}

	task := &roomsmodels.RoomTaskModel{
		TaskID:      uuid.NewString(),
		RoomID:      roomID,
		Title:       title,
		Description: description,
		ExternalKey: externalKey,
		Status:      "PENDING",
		IsActive:    false,
	}

	createdTask, err := s.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	s.expiryService.TouchActivity(roomID)

	return createdTask, nil
}

func (s *roomsTaskService) ListTasks(roomID, userID string) ([]*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	return s.taskRepo.FindByRoomID(roomID)
}

func (s *roomsTaskService) GetTask(roomID, taskID, userID string) (*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	return s.taskRepo.FindByID(roomID, taskID)
}

func (s *roomsTaskService) UpdateTask(roomID, taskID, userID string, input UpdateTaskInput) (*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureActiveRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.FindByID(roomID, taskID)
	if err != nil {
		return nil, err
	}

	if hasVoteStateInput(input) {
		task, err = s.applyVoteStateUpdate(roomID, task, userID, input)
		if err != nil {
			return nil, err
		}
	}

	changed, err := applyTaskMetadata(task, input)
	if err != nil {
		return nil, err
	}
	if !changed {
		return task, nil
	}

	updatedTask, err := s.taskRepo.Update(roomID, task)
	if err != nil {
		return nil, err
	}

	s.expiryService.TouchActivity(roomID)

	return updatedTask, nil
}

func (s *roomsTaskService) DeleteTask(roomID, taskID, userID string) error {
	if _, err := s.ensureActiveRoomAdmin(roomID, userID); err != nil {
		return err
	}

	if err := s.taskRepo.Delete(roomID, taskID); err != nil {
		return err
	}

	s.expiryService.TouchActivity(roomID)

	return nil
}

func (s *roomsTaskService) applyVoteStateUpdate(
	roomID string,
	task *roomsmodels.RoomTaskModel,
	userID string,
	input UpdateTaskInput,
) (*roomsmodels.RoomTaskModel, error) {
	status := ""
	if input.Status != nil {
		status = strings.TrimSpace(*input.Status)
	}

	activateRequested := (input.IsActive != nil && *input.IsActive) || status == "VOTING"
	deactivateRequested := input.IsActive != nil && !*input.IsActive
	finalizeRequested := input.FinalEstimateValue != nil || status == "ESTIMATED"

	if activateRequested && (finalizeRequested || status == "SKIPPED" || status == "PENDING" || deactivateRequested) {
		return nil, fmt.Errorf("%w: conflicting task state update", apperrors.ErrBadRequest)
	}
	if status == "VOTING" && input.IsActive != nil && !*input.IsActive {
		return nil, fmt.Errorf("%w: active task must remain active while voting", apperrors.ErrBadRequest)
	}

	switch {
	case activateRequested:
		updatedTask, _, _, err := s.voteService.SetCurrentTask(roomID, task.TaskID, userID, nil)
		if err != nil {
			return nil, err
		}
		return updatedTask, nil
	case finalizeRequested:
		value := ""
		if input.FinalEstimateValue != nil {
			value = strings.TrimSpace(*input.FinalEstimateValue)
		}
		if value == "" {
			if task.FinalEstimateValue == nil || strings.TrimSpace(*task.FinalEstimateValue) == "" {
				return nil, fmt.Errorf("%w: final estimate value is required", apperrors.ErrBadRequest)
			}
			value = strings.TrimSpace(*task.FinalEstimateValue)
		}
		return s.voteService.FinalizeTask(roomID, task.TaskID, userID, value)
	case status == "SKIPPED":
		task.Status = "SKIPPED"
		task.IsActive = false
		task.FinalEstimateValue = nil
		updatedTask, err := s.taskRepo.Update(roomID, task)
		if err != nil {
			return nil, err
		}
		s.expiryService.TouchActivity(roomID)
		return updatedTask, nil
	case status == "PENDING" || deactivateRequested:
		task.Status = "PENDING"
		task.IsActive = false
		updatedTask, err := s.taskRepo.Update(roomID, task)
		if err != nil {
			return nil, err
		}
		s.expiryService.TouchActivity(roomID)
		return updatedTask, nil
	default:
		return task, nil
	}
}

func hasVoteStateInput(input UpdateTaskInput) bool {
	return input.Status != nil || input.IsActive != nil || input.FinalEstimateValue != nil
}

func applyTaskMetadata(task *roomsmodels.RoomTaskModel, input UpdateTaskInput) (bool, error) {
	changed := false

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return false, apperrors.ErrBadRequest
		}
		if task.Title != title {
			task.Title = title
			changed = true
		}
	}

	if input.Description != nil {
		trimmed := strings.TrimSpace(*input.Description)
		if trimmed == "" {
			if task.Description != nil {
				task.Description = nil
				changed = true
			}
		} else if task.Description == nil || *task.Description != *input.Description {
			task.Description = input.Description
			changed = true
		}
	}

	if input.ExternalKey != nil {
		trimmed := strings.TrimSpace(*input.ExternalKey)
		if trimmed == "" {
			if task.ExternalKey != nil {
				task.ExternalKey = nil
				changed = true
			}
		} else if task.ExternalKey == nil || *task.ExternalKey != *input.ExternalKey {
			task.ExternalKey = input.ExternalKey
			changed = true
		}
	}

	return changed, nil
}

func (s *roomsTaskService) ensureRoomAdmin(roomID, userID string) (*roomsmodels.RoomsModel, error) {
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

func (s *roomsTaskService) ensureActiveRoomAdmin(roomID, userID string) (*roomsmodels.RoomsModel, error) {
	room, err := s.ensureRoomAdmin(roomID, userID)
	if err != nil {
		return nil, err
	}
	if room.Status != "ACTIVE" {
		return nil, apperrors.ErrForbidden
	}

	return room, nil
}
