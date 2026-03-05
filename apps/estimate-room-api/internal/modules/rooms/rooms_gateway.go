package rooms

import (
	"encoding/json"
	"errors"
	"strings"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type roomsGateway struct {
	wsService       *ws.Service
	participantRepo roomsrepositories.RoomParticipantRepository
}

func NewRoomsGateway(
	wsService *ws.Service,
	participantRepo roomsrepositories.RoomParticipantRepository,
) *roomsGateway {
	return &roomsGateway{
		wsService:       wsService,
		participantRepo: participantRepo,
	}
}

const (
	RoomsJoin           = "ROOMS_JOIN"
	RoomsTaskSetCurrent = "ROOMS_TASK_SET_CURRENT"
	RoomsVoteCast       = "ROOMS_VOTE_CAST"
	RoomsVoteReveal     = "ROOMS_VOTE_REVEAL"
	RoomsRoundNext      = "ROOMS_ROUND_NEXT"

	RoomsParticipantJoined = "ROOMS_PARTICIPANT_JOINED"
	RoomsParticipantLeft   = "ROOMS_PARTICIPANT_LEFT"
)

type roomJoinPayload struct {
	RoomID string `json:"roomId"`
}

type roomPresencePayload struct {
	ParticipantID string                          `json:"participantId,omitempty"`
	UserID        *string                         `json:"userId,omitempty"`
	GuestName     *string                         `json:"guestName,omitempty"`
	Role          roomsmodels.RoomParticipantRole `json:"role,omitempty"`
}

func (g *roomsGateway) handleRoomJoin(client ws.ClientInfo, event ws.Event) {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID == "" && len(event.Payload) > 0 {
		payload := roomJoinPayload{}
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			logger.L().Error("room join payload parse failed", "err", err, "user_id", client.UserID, "conn_id", client.ConnID)
			return
		}
		roomID = strings.TrimSpace(payload.RoomID)
	}

	if roomID == "" {
		logger.L().Warn("room join ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	participant, err := g.resolveParticipant(client, roomID)
	if err != nil {
		logJoinDenied(client, roomID, err)
		return
	}

	if err := g.wsService.SetParticipantID(client.ConnID, participant.RoomParticipantID); err != nil {
		logger.L().Error("failed to bind ws participant", "err", err, "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	joinResult, err := g.wsService.JoinRoom(client.ConnID, roomID)
	if err != nil {
		logger.L().Error("room join failed", "err", err, "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	if joinResult.Joined {
		if err := g.broadcastPresence(roomID, RoomsParticipantJoined, roomPresencePayload{
			ParticipantID: participant.RoomParticipantID,
			UserID:        participant.UserID,
			GuestName:     participant.GuestName,
			Role:          participant.Role,
		}); err != nil {
			logger.L().Error("failed to broadcast participant joined", "err", err, "room_id", roomID)
		}
	}

	logger.L().Info("room join accepted", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleTaskSetCurrent(client ws.ClientInfo, event ws.Event) {
	logger.L().Info("task set current received", "room_id", event.RoomID, "user_id", client.UserID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleVoteCast(client ws.ClientInfo, event ws.Event) {
	logger.L().Info("vote cast received", "room_id", event.RoomID, "user_id", client.UserID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleVoteReveal(client ws.ClientInfo, event ws.Event) {
	logger.L().Info("vote reveal received", "room_id", event.RoomID, "user_id", client.UserID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleRoundNext(client ws.ClientInfo, event ws.Event) {
	logger.L().Info("round next received", "room_id", event.RoomID, "user_id", client.UserID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleDisconnect(info ws.DisconnectInfo) {
	roomID := strings.TrimSpace(info.RoomID)
	if roomID == "" || !info.PresenceLeft {
		return
	}

	participantID := strings.TrimSpace(info.Client.ParticipantID)
	payload := roomPresencePayload{
		ParticipantID: participantID,
	}

	if info.Client.UserID != "" {
		userID := info.Client.UserID
		payload.UserID = &userID
	}

	if err := g.broadcastPresence(roomID, RoomsParticipantLeft, payload); err != nil {
		logger.L().Error("failed to broadcast participant left", "err", err, "room_id", roomID)
	}
}

func (g *roomsGateway) resolveParticipant(client ws.ClientInfo, roomID string) (*roomsmodels.RoomParticipantModel, error) {
	switch client.IdentityType {
	case ws.IdentityTypeUser:
		userID := strings.TrimSpace(client.UserID)
		if userID == "" {
			return nil, apperrors.ErrUnauthorized
		}
		return g.participantRepo.FindActiveByUserID(roomID, userID)
	case ws.IdentityTypeGuest:
		participantID := strings.TrimSpace(client.ParticipantID)
		if participantID == "" {
			return nil, apperrors.ErrUnauthorized
		}
		participant, err := g.participantRepo.FindActiveByID(roomID, participantID)
		if err != nil {
			return nil, err
		}
		if participant.Role != roomsmodels.RoomParticipantRoleGuest {
			return nil, apperrors.ErrForbidden
		}
		return participant, nil
	default:
		return nil, apperrors.ErrUnauthorized
	}
}

func (g *roomsGateway) broadcastPresence(roomID, eventType string, payload roomPresencePayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    eventType,
		RoomID:  roomID,
		Payload: data,
	})
}

func logJoinDenied(client ws.ClientInfo, roomID string, err error) {
	switch {
	case errors.Is(err, apperrors.ErrForbidden), errors.Is(err, apperrors.ErrUnauthorized), errors.Is(err, apperrors.ErrNotFound):
		logger.L().Warn("room join denied", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID, "reason", err.Error())
	default:
		logger.L().Error("room join failed", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID, "err", err)
	}
}
