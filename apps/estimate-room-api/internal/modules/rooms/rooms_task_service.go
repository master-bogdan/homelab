package rooms

import (
	"errors"
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
	participantRepo roomsrepositories.RoomParticipantRepository
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
	FinalEstimateValue *string
}

func NewRoomsTaskService(
	roomsRepo roomsrepositories.RoomsRepository,
	taskRepo roomsrepositories.RoomTaskRepository,
	participantRepo roomsrepositories.RoomParticipantRepository,
) RoomsTaskService {
	return &roomsTaskService{
		roomsRepo:       roomsRepo,
		taskRepo:        taskRepo,
		participantRepo: participantRepo,
		logger:          logger.L().With(slog.String("service", "rooms-tasks")),
	}
}

func (s *roomsTaskService) CreateTask(roomID, userID string, input CreateTaskInput) (*roomsmodels.RoomTaskModel, error) {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
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
	}

	return s.taskRepo.Create(task)
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
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return nil, err
	}

	task, err := s.taskRepo.FindByID(roomID, taskID)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, apperrors.ErrBadRequest
		}
		task.Title = title
	}

	if input.Description != nil {
		if *input.Description == "" {
			task.Description = nil
		} else {
			task.Description = input.Description
		}
	}

	if input.ExternalKey != nil {
		if *input.ExternalKey == "" {
			task.ExternalKey = nil
		} else {
			task.ExternalKey = input.ExternalKey
		}
	}

	if input.Status != nil {
		task.Status = strings.TrimSpace(*input.Status)
	}

	if input.FinalEstimateValue != nil {
		if *input.FinalEstimateValue == "" {
			task.FinalEstimateValue = nil
		} else {
			task.FinalEstimateValue = input.FinalEstimateValue
		}
	}

	return s.taskRepo.Update(roomID, task)
}

func (s *roomsTaskService) DeleteTask(roomID, taskID, userID string) error {
	if _, err := s.ensureRoomAdmin(roomID, userID); err != nil {
		return err
	}

	return s.taskRepo.Delete(roomID, taskID)
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
