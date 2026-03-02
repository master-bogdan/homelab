package rooms

import (
	"strings"

	"github.com/google/uuid"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

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

func (s *roomsService) CreateTask(roomID string, input CreateTaskInput) (*roomsmodels.RoomTaskModel, error) {
	if err := s.ensureRoomExists(roomID); err != nil {
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

	return s.roomsRepo.CreateTask(task)
}

func (s *roomsService) ListTasks(roomID string) ([]*roomsmodels.RoomTaskModel, error) {
	if err := s.ensureRoomExists(roomID); err != nil {
		return nil, err
	}

	return s.roomsRepo.FindTasksByRoomID(roomID)
}

func (s *roomsService) GetTask(roomID, taskID string) (*roomsmodels.RoomTaskModel, error) {
	return s.roomsRepo.FindTaskByID(roomID, taskID)
}

func (s *roomsService) UpdateTask(roomID, taskID string, input UpdateTaskInput) (*roomsmodels.RoomTaskModel, error) {
	task, err := s.roomsRepo.FindTaskByID(roomID, taskID)
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

	return s.roomsRepo.UpdateTask(roomID, task)
}

func (s *roomsService) DeleteTask(roomID, taskID string) error {
	return s.roomsRepo.DeleteTask(roomID, taskID)
}

func (s *roomsService) ensureRoomExists(roomID string) error {
	exists, err := s.roomsRepo.RoomExists(roomID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.ErrNotFound
	}

	return nil
}
