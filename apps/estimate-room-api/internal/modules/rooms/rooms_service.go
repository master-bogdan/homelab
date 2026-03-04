package rooms

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsService interface {
	CreateRoom(model roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error)
	GetRoom(roomID string) (*roomsmodels.RoomsModel, error)
	ValidateUserRoomAccess(roomID, userID string) error
	UpdateRoom(roomID, userID string, input UpdateRoomInput) (*roomsmodels.RoomsModel, error)
}

type roomsService struct {
	roomsRepo       roomsrepositories.RoomsRepository
	participantRepo roomsrepositories.RoomParticipantRepository
	logger          *slog.Logger
}

func NewRoomsService(
	roomsRepo roomsrepositories.RoomsRepository,
	participantRepo roomsrepositories.RoomParticipantRepository,
) RoomsService {
	return &roomsService{
		roomsRepo:       roomsRepo,
		participantRepo: participantRepo,
		logger:          logger.L().With(slog.String("service", "rooms")),
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

	room, err := s.roomsRepo.Create(&model)
	if err != nil {
		return nil, err
	}

	_, err = s.participantRepo.Create(&roomsmodels.RoomParticipantModel{
		RoomParticipantID: uuid.NewString(),
		RoomID:            room.RoomID,
		UserID:            &room.AdminUserID,
		Role:              roomsmodels.RoomParticipantRoleAdmin,
	})
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomsService) GetRoom(roomID string) (*roomsmodels.RoomsModel, error) {
	return s.roomsRepo.FindByID(roomID)
}

func (s *roomsService) ValidateUserRoomAccess(roomID, userID string) error {
	participant, err := s.participantRepo.FindActiveByUserID(roomID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrForbidden
		}
		return err
	}

	if !participant.Role.IsValid() {
		return apperrors.ErrForbidden
	}

	return nil
}

type UpdateRoomInput struct {
	Name   *string
	Status *string
}

func (s *roomsService) UpdateRoom(roomID, userID string, input UpdateRoomInput) (*roomsmodels.RoomsModel, error) {
	room, err := s.ensureRoomAdmin(roomID, userID)
	if err != nil {
		return nil, err
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
		Name:   input.Name,
		Status: input.Status,
	})
}

func roomPatchIsNoop(room *roomsmodels.RoomsModel, input UpdateRoomInput) bool {
	if input.Name != nil && room.Name != *input.Name {
		return false
	}
	if input.Status != nil && room.Status != *input.Status {
		return false
	}
	return true
}

func (s *roomsService) ensureRoomAdmin(roomID, userID string) (*roomsmodels.RoomsModel, error) {
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
