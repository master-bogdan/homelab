package rooms

import (
	"encoding/json"
	"errors"

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
	EventRoomJoin    = "ROOM_JOIN"
	EventRoomLeave   = "ROOM_LEAVE"
	EventRoomMessage = "ROOM_MESSAGE"
)

func (g *roomsGateway) OnEvent(client ws.ClientInfo, event ws.Event) {
	logger.L().Info("Event received", "type", event.Type, "user_id", client.UserID, "conn_id", client.ConnID)

	var payload map[string]any
	if len(event.Payload) > 0 {
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			logger.L().Error("Invalid event payload", "err", err)
			return
		}
	}
}

func (g *roomsGateway) SendToRoom(channelID string, data any) error {
	if channelID == "" {
		return errors.New("channelID is required")
	}

	payload, err := json.Marshal(map[string]any{
		"roomId": channelID,
		"data":   data,
	})
	if err != nil {
		return err
	}

	event := ws.Event{
		Type:    EventRoomMessage,
		Payload: payload,
	}

	return g.wsService.Broadcast(event)
}
