package rooms

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (g *roomsGateway) HandleConnection(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "roomID")
	if channelID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	clientID := r.URL.Query().Get("clientID")
	if clientID == "" {
		anonID, err := newClientID()
		if err != nil {
			http.Error(w, "failed to generate clientID", http.StatusInternalServerError)
			return
		}
		clientID = anonID
	}

	g.wsManager.HandleWS(w, r, channelID, clientID, g)
}

func (g *roomsGateway) OnConnect(client ws.ClientInfo) {
	log.Printf("Client %s connected to channel %s", client.ClientID, client.ChannelID)
}

func (g *roomsGateway) OnDisconnect(client ws.ClientInfo) {
	log.Printf("Client %s disconnected from channel %s", client.ClientID, client.ChannelID)
}

func (g *roomsGateway) OnMessage(client ws.ClientInfo, message []byte) {
	log.Printf("Message from client %s in channel %s: %s", client.ClientID, client.ChannelID, string(message))

	var msg map[string]any
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Invalid JSON: %v", err)
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
		log.Printf("Broadcast error: %v", err)
		return
	}
}

func (g *roomsGateway) OnError(client ws.ClientInfo, err error) {
	log.Printf("Error for client %s in channel %s: %v", client.ClientID, client.ChannelID, err)
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
