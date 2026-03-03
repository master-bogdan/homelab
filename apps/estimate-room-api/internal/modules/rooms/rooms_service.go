package rooms

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"strings"
	"time"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsService interface {
	CreateRoom(model roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error)
	GetRoom(roomID string) (*roomsmodels.RoomsModel, error)
	UpdateRoom(roomID, userID string, input UpdateRoomInput) (*roomsmodels.RoomsModel, error)
	CreateTask(roomID string, input CreateTaskInput) (*roomsmodels.RoomTaskModel, error)
	ListTasks(roomID string) ([]*roomsmodels.RoomTaskModel, error)
	GetTask(roomID, taskID string) (*roomsmodels.RoomTaskModel, error)
	UpdateTask(roomID, taskID string, input UpdateTaskInput) (*roomsmodels.RoomTaskModel, error)
	DeleteTask(roomID, taskID string) error
}

type roomsService struct {
	roomsRepo roomsrepositories.RoomsRepository
	logger    *slog.Logger
}

func NewRoomsService(roomsRepo roomsrepositories.RoomsRepository) RoomsService {
	return &roomsService{
		roomsRepo: roomsRepo,
		logger:    logger.L().With(slog.String("service", "rooms")),
	}
}

func (s *roomsService) CreateRoom(model roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error) {
	if model.DeckID == "" {
		model.DeckID = roomsmodels.DeckIDFibonacci
	}
	if !model.DeckID.IsValid() {
		return nil, errors.New("invalid deck id")
	}

	timestamp := time.Now().String()
	code := base64.StdEncoding.EncodeToString([]byte(model.Name + timestamp))

	model.Code = code

	return s.roomsRepo.Create(&model)
}

func (s *roomsService) GetRoom(roomID string) (*roomsmodels.RoomsModel, error) {
	return s.roomsRepo.FindByID(roomID)
}

type UpdateRoomInput struct {
	Name              *string
	Status            *string
	AllowGuests       *bool
	AllowSpectators   *bool
	RoundTimerSeconds *int
}

func (s *roomsService) UpdateRoom(roomID, userID string, input UpdateRoomInput) (*roomsmodels.RoomsModel, error) {
	room, err := s.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	if room.AdminUserID != userID {
		return nil, apperrors.ErrForbidden
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, apperrors.ErrBadRequest
		}
		input.Name = &name
	}

	if roomPatchIsNoop(room, input) {
		return room, nil
	}

	return s.roomsRepo.Update(roomID, roomsrepositories.UpdateRoomFields{
		Name:              input.Name,
		Status:            input.Status,
		AllowGuests:       input.AllowGuests,
		AllowSpectators:   input.AllowSpectators,
		RoundTimerSeconds: input.RoundTimerSeconds,
	})
}

func roomPatchIsNoop(room *roomsmodels.RoomsModel, input UpdateRoomInput) bool {
	if input.Name != nil && room.Name != *input.Name {
		return false
	}
	if input.Status != nil && room.Status != *input.Status {
		return false
	}
	if input.AllowGuests != nil && room.AllowGuests != *input.AllowGuests {
		return false
	}
	if input.AllowSpectators != nil && room.AllowSpectators != *input.AllowSpectators {
		return false
	}
	if input.RoundTimerSeconds != nil && room.RoundTimerSeconds != *input.RoundTimerSeconds {
		return false
	}

	return true
}
