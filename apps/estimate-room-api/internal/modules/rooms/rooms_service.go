package rooms

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsService interface {
	CreateRoom(model roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error)
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
