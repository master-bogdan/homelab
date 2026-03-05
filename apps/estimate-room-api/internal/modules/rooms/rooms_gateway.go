package rooms

import (
	"encoding/json"
	"strings"

	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type roomsGateway struct {
	wsService *ws.Service
}

func NewRoomsGateway(
	wsService *ws.Service,
) *roomsGateway {
	return &roomsGateway{
		wsService: wsService,
	}
}

const (
	RoomsJoin           = "ROOMS_JOIN"
	RoomsTaskSetCurrent = "ROOMS_TASK_SET_CURRENT"
	RoomsVoteCast       = "ROOMS_VOTE_CAST"
	RoomsVoteReveal     = "ROOMS_VOTE_REVEAL"
	RoomsRoundNext      = "ROOMS_ROUND_NEXT"
)

type roomJoinPayload struct {
	RoomID string `json:"roomId"`
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

	logger.L().Info("room join received", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID)
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
