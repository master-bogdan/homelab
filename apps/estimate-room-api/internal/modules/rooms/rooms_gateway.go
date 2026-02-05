package rooms

import (
	"encoding/json"
	"errors"

	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
)

type roomsGateway struct {
	wsManager *ws.Manager
}

func NewRoomsGateway(
	wsManager *ws.Manager,
) *roomsGateway {
	return &roomsGateway{
		wsManager: wsManager,
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

// SendToRoom publishes a server message to all clients in the room.
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

	return g.wsManager.Broadcast(event)
}
