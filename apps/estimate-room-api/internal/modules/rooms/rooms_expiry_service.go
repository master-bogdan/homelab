package rooms

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

const (
	roomExpiryAfter         = 30 * time.Minute
	roomExpirySweepInterval = 1 * time.Minute
)

type RoomsExpiryService interface {
	TouchActivity(roomID string)
	ExpireInactiveRooms(cutoff time.Time) ([]*roomsmodels.RoomsModel, error)
	Start(ctx context.Context)
}

type roomsExpiryService struct {
	roomsRepo roomsrepositories.RoomsRepository
	wsService *ws.Service
	logger    *slog.Logger
}

type roomExpiredPayload struct {
	RoomID string `json:"roomId"`
	Status string `json:"status"`
}

func NewRoomsExpiryService(
	roomsRepo roomsrepositories.RoomsRepository,
	wsService *ws.Service,
) RoomsExpiryService {
	return &roomsExpiryService{
		roomsRepo: roomsRepo,
		wsService: wsService,
		logger:    logger.L().With(slog.String("service", "rooms-expiry")),
	}
}

// Activity is any successful collaborative interaction that changes room presence
// or mutable room state: websocket room joins/leaves, invite joins, room updates,
// task changes, and the vote lifecycle.
func (s *roomsExpiryService) TouchActivity(roomID string) {
	if strings.TrimSpace(roomID) == "" {
		return
	}

	if err := s.roomsRepo.TouchActivity(roomID); err != nil {
		s.logger.Error("failed to touch room activity", "room_id", roomID, "err", err)
	}
}

func (s *roomsExpiryService) ExpireInactiveRooms(cutoff time.Time) ([]*roomsmodels.RoomsModel, error) {
	expiredRooms, err := s.roomsRepo.ExpireInactiveRooms(cutoff)
	if err != nil {
		return nil, err
	}

	for _, room := range expiredRooms {
		if room == nil {
			continue
		}
		if err := s.broadcastRoomExpired(room); err != nil {
			s.logger.Error("failed to broadcast room expired", "room_id", room.RoomID, "err", err)
		}
	}

	return expiredRooms, nil
}

func (s *roomsExpiryService) Start(ctx context.Context) {
	if ctx == nil {
		return
	}

	go func() {
		s.runExpirySweep()

		ticker := time.NewTicker(roomExpirySweepInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.runExpirySweep()
			}
		}
	}()
}

func (s *roomsExpiryService) runExpirySweep() {
	cutoff := time.Now().Add(-roomExpiryAfter)
	expiredRooms, err := s.ExpireInactiveRooms(cutoff)
	if err != nil {
		s.logger.Error("failed to expire inactive rooms", "err", err)
		return
	}

	if len(expiredRooms) > 0 {
		s.logger.Info("expired inactive rooms", "count", len(expiredRooms), "cutoff", cutoff)
	}
}

func (s *roomsExpiryService) broadcastRoomExpired(room *roomsmodels.RoomsModel) error {
	if room == nil || s.wsService == nil {
		return nil
	}

	data, err := json.Marshal(roomExpiredPayload{
		RoomID: room.RoomID,
		Status: room.Status,
	})
	if err != nil {
		return err
	}

	return s.wsService.Broadcast(ws.Event{
		Type:    RoomsExpired,
		RoomID:  room.RoomID,
		Payload: data,
	})
}
