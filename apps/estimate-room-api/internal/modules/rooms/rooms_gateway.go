package rooms

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
)

type roomsGateway struct {
	wsManager *ws.Manager
}

func NewRoomsGateway(
	wsManager *ws.Manager,
) ws.Gateway {
	return &roomsGateway{
		wsManager: wsManager,
	}
}

// HandleConnection godoc
// @Summary Room WebSocket connection
// @Description Upgrades the HTTP request to a WebSocket for the given room.
// @Tags rooms
// @Param roomID path string true "Room ID"
// @Param clientID query string false "Client ID"
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/{roomID}/ws [get]
func (g *roomsGateway) HandleConnection(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "roomID")
	if channelID == "" {
		utils.WriteResponseError(w, http.StatusBadRequest, "roomID is required")
		return
	}

	clientID := r.URL.Query().Get("clientID")
	if clientID == "" {
		anonID, err := newClientID()
		if err != nil {
			utils.WriteResponseError(w, http.StatusInternalServerError, "failed to generate clientID")
			return
		}
		clientID = anonID
	}

	g.wsManager.HandleWS(w, r, channelID, clientID, g)
}

func (g *roomsGateway) OnConnect(client ws.ClientInfo) {
	logger.L().Info("Client connected", "client_id", client.ClientID, "channel_id", client.ChannelID)
}

func (g *roomsGateway) OnDisconnect(client ws.ClientInfo) {
	logger.L().Info("Client disconnected", "client_id", client.ClientID, "channel_id", client.ChannelID)
}

func (g *roomsGateway) OnMessage(client ws.ClientInfo, message []byte) {
	logger.L().Info("Message received", "client_id", client.ClientID, "channel_id", client.ChannelID, "message", string(message))

	var msg map[string]any
	if err := json.Unmarshal(message, &msg); err != nil {
		logger.L().Error("Invalid JSON", "err", err)
		return
	}

	response := map[string]any{
		"type":      "message",
		"channelID": client.ChannelID,
		"data":      msg,
	}

	// Broadcast via Redis to all servers
	err := g.wsManager.Broadcast(response)
	if err != nil {
		logger.L().Error("Broadcast error", "err", err)
		return
	}
}

func (g *roomsGateway) OnError(client ws.ClientInfo, err error) {
	logger.L().Error("Client error", "client_id", client.ClientID, "channel_id", client.ChannelID, "err", err)
}

// SendToRoom publishes a server message to all clients in the room.
func (g *roomsGateway) SendToRoom(channelID string, data any) error {
	if channelID == "" {
		return errors.New("channelID is required")
	}

	response := map[string]any{
		"type":      "server-message",
		"channelID": channelID,
		"data":      data,
	}

	return g.wsManager.Broadcast(response)
}

func newClientID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw[:]), nil
}
