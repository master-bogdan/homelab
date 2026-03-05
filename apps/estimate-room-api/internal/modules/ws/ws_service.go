package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

const (
	pingInterval = 30 * time.Second
	pingTimeout  = 10 * time.Second
)

type Client struct {
	Conn          *websocket.Conn
	ConnID        string
	IdentityType  IdentityType
	IdentityID    string
	UserID        string
	ParticipantID string
	RoomID        string
	Send          chan []byte
}

type Service struct {
	clients         map[*Client]bool
	connClients     map[string]*Client
	identityClients map[string]map[*Client]bool
	roomClients     map[string]map[*Client]bool
	roomPresence    map[string]map[string]int
	register        chan *Client
	unregister      chan *Client
	mu              sync.RWMutex
	server          PubSub
	channel         string
	subscriptions   map[string][]EventHandler
	disconnects     []DisconnectHandler
}

func NewService(server PubSub, channel string) *Service {
	s := &Service{
		clients:         make(map[*Client]bool),
		connClients:     make(map[string]*Client),
		identityClients: make(map[string]map[*Client]bool),
		roomClients:     make(map[string]map[*Client]bool),
		roomPresence:    make(map[string]map[string]int),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		server:          server,
		channel:         channel,
		subscriptions:   make(map[string][]EventHandler),
		disconnects:     make([]DisconnectHandler, 0),
	}

	go s.run()

	if server != nil {
		server.Subscribe(channel, func(data []byte) {
			s.broadcastRaw(data)
		})
	}

	return s
}

func clientInfo(client *Client) ClientInfo {
	if client == nil {
		return ClientInfo{}
	}
	return ClientInfo{
		ConnID:        client.ConnID,
		IdentityType:  client.IdentityType,
		IdentityID:    client.IdentityID,
		UserID:        client.UserID,
		ParticipantID: client.ParticipantID,
	}
}

func (s *Service) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.connClients[client.ConnID] = client

			if s.identityClients[client.IdentityID] == nil {
				s.identityClients[client.IdentityID] = make(map[*Client]bool)
			}
			s.identityClients[client.IdentityID][client] = true
			s.mu.Unlock()

			s.logConnect(clientInfo(client))

		case client := <-s.unregister:
			var disconnectInfo *DisconnectInfo

			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				roomID, presenceLeft := s.removeClientFromRoomLocked(client)

				delete(s.clients, client)
				delete(s.connClients, client.ConnID)
				delete(s.identityClients[client.IdentityID], client)
				if len(s.identityClients[client.IdentityID]) == 0 {
					delete(s.identityClients, client.IdentityID)
				}

				if roomID != "" {
					info := DisconnectInfo{
						Client:       clientInfo(client),
						RoomID:       roomID,
						PresenceLeft: presenceLeft,
					}
					disconnectInfo = &info
				}

				close(client.Send)
			}
			s.mu.Unlock()

			s.logDisconnect(clientInfo(client))
			if disconnectInfo != nil {
				s.dispatchDisconnect(*disconnectInfo)
			}
		}
	}
}

func (s *Service) Subscribe(eventType string, handler EventHandler) {
	if strings.TrimSpace(eventType) == "" || handler == nil {
		return
	}
	s.mu.Lock()
	s.subscriptions[eventType] = append(s.subscriptions[eventType], handler)
	s.mu.Unlock()
}

func (s *Service) SubscribeDisconnect(handler DisconnectHandler) {
	if handler == nil {
		return
	}

	s.mu.Lock()
	s.disconnects = append(s.disconnects, handler)
	s.mu.Unlock()
}

type JoinRoomResult struct {
	RoomID         string
	Joined         bool
	PreviousRoomID string
	PreviousLeft   bool
}

func (s *Service) JoinRoom(connID, roomID string) (JoinRoomResult, error) {
	trimmedConnID := strings.TrimSpace(connID)
	trimmedRoomID := strings.TrimSpace(roomID)
	if trimmedConnID == "" || trimmedRoomID == "" {
		return JoinRoomResult{}, errors.New("connID and roomID are required")
	}

	result := JoinRoomResult{RoomID: trimmedRoomID}

	s.mu.Lock()
	defer s.mu.Unlock()

	client, ok := s.connClients[trimmedConnID]
	if !ok || client == nil {
		return JoinRoomResult{}, errors.New("connection not found")
	}

	if client.RoomID == trimmedRoomID {
		return result, nil
	}

	if client.RoomID != "" {
		prevRoomID, prevLeft := s.removeClientFromRoomLocked(client)
		result.PreviousRoomID = prevRoomID
		result.PreviousLeft = prevLeft
	}

	if s.roomClients[trimmedRoomID] == nil {
		s.roomClients[trimmedRoomID] = make(map[*Client]bool)
	}
	s.roomClients[trimmedRoomID][client] = true

	if s.roomPresence[trimmedRoomID] == nil {
		s.roomPresence[trimmedRoomID] = make(map[string]int)
	}
	s.roomPresence[trimmedRoomID][client.IdentityID]++
	if s.roomPresence[trimmedRoomID][client.IdentityID] == 1 {
		result.Joined = true
	}

	client.RoomID = trimmedRoomID

	return result, nil
}

func (s *Service) SetParticipantID(connID, participantID string) error {
	trimmedConnID := strings.TrimSpace(connID)
	trimmedParticipantID := strings.TrimSpace(participantID)
	if trimmedConnID == "" || trimmedParticipantID == "" {
		return errors.New("connID and participantID are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	client, ok := s.connClients[trimmedConnID]
	if !ok || client == nil {
		return errors.New("connection not found")
	}

	client.ParticipantID = trimmedParticipantID
	return nil
}

func (s *Service) Connect(w http.ResponseWriter, r *http.Request, identity ConnectIdentity) {
	identityID, ok := resolveIdentityID(identity)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		logger.L().Error("WebSocket accept error", "err", err)
		return
	}

	connID := uuid.NewString()
	client := &Client{
		Conn:          conn,
		ConnID:        connID,
		IdentityType:  identity.Type,
		IdentityID:    identityID,
		UserID:        strings.TrimSpace(identity.UserID),
		ParticipantID: strings.TrimSpace(identity.ParticipantID),
		Send:          make(chan []byte, 256),
	}

	s.disconnectIdentityConnections(identityID)

	s.register <- client

	go s.writeHandler(client)
	s.sendHello(client)
	s.readHandler(client)
}

func (s *Service) readHandler(client *Client) {
	defer func() {
		s.unregister <- client
		client.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	ctx := context.Background()
	for {
		_, message, err := client.Conn.Read(ctx)
		if err != nil {
			s.logError(clientInfo(client), err)
			break
		}
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			s.logError(clientInfo(client), err)
			continue
		}
		s.normalizeIncomingEvent(client, &event)
		if strings.TrimSpace(event.Type) == "" {
			s.logError(clientInfo(client), errors.New("missing event type"))
			continue
		}
		s.dispatchEvent(clientInfo(client), event)
	}
}

func (s *Service) writeHandler(client *Client) {
	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}
			ctx := context.Background()
			err := client.Conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				s.logError(clientInfo(client), err)
				return
			}
		case <-pingTicker.C:
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
			err := client.Conn.Ping(ctx)
			cancel()
			if err != nil {
				s.logError(clientInfo(client), err)
				client.Conn.Close(websocket.StatusNormalClosure, "ping timeout")
				return
			}
		}
	}
}

func (s *Service) broadcastRaw(data []byte) {
	event := Event{}
	if err := json.Unmarshal(data, &event); err == nil {
		roomID := strings.TrimSpace(event.RoomID)
		if roomID != "" {
			s.broadcastRoomRaw(roomID, data)
			return
		}
	}

	s.mu.RLock()
	clients := make([]*Client, 0, len(s.clients))
	for client := range s.clients {
		clients = append(clients, client)
	}
	s.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			s.unregister <- client
		}
	}
}

func (s *Service) broadcastRoomRaw(roomID string, data []byte) {
	s.mu.RLock()
	roomMembers := s.roomClients[roomID]
	clients := make([]*Client, 0, len(roomMembers))
	for client := range roomMembers {
		clients = append(clients, client)
	}
	s.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			s.unregister <- client
		}
	}
}

func (s *Service) Broadcast(message any) error {
	if s.server == nil {
		return errors.New("ws server is not initialized")
	}

	switch v := message.(type) {
	case Event:
		s.normalizeOutgoingEvent(&v)
		return s.server.Publish(s.channel, v)
	case *Event:
		if v == nil {
			return errors.New("ws event is nil")
		}
		s.normalizeOutgoingEvent(v)
		return s.server.Publish(s.channel, v)
	}

	return s.server.Publish(s.channel, message)
}

func (s *Service) disconnectIdentityConnections(identityID string) {
	s.mu.RLock()
	existing := make([]*Client, 0, len(s.identityClients[identityID]))
	for client := range s.identityClients[identityID] {
		existing = append(existing, client)
	}
	s.mu.RUnlock()

	for _, client := range existing {
		_ = client.Conn.Close(websocket.StatusPolicyViolation, "another connection opened")
	}
}

func (s *Service) dispatchEvent(info ClientInfo, event Event) {
	s.mu.RLock()
	handlers := append([]EventHandler(nil), s.subscriptions[event.Type]...)
	s.mu.RUnlock()

	for _, handler := range handlers {
		handler(info, event)
	}
}

func (s *Service) dispatchDisconnect(info DisconnectInfo) {
	s.mu.RLock()
	handlers := append([]DisconnectHandler(nil), s.disconnects...)
	s.mu.RUnlock()

	for _, handler := range handlers {
		handler(info)
	}
}

func (s *Service) logConnect(info ClientInfo) {
	logger.L().Info(
		"ws connected",
		"conn_id", info.ConnID,
		"identity_type", info.IdentityType,
		"identity_id", info.IdentityID,
		"user_id", info.UserID,
		"participant_id", info.ParticipantID,
	)
}

func (s *Service) logDisconnect(info ClientInfo) {
	logger.L().Info(
		"ws disconnected",
		"conn_id", info.ConnID,
		"identity_type", info.IdentityType,
		"identity_id", info.IdentityID,
		"user_id", info.UserID,
		"participant_id", info.ParticipantID,
	)
}

func (s *Service) logError(info ClientInfo, err error) {
	logger.L().Error(
		"ws error",
		"conn_id", info.ConnID,
		"identity_type", info.IdentityType,
		"identity_id", info.IdentityID,
		"user_id", info.UserID,
		"participant_id", info.ParticipantID,
		"err", err,
	)
}

func (s *Service) sendHello(client *Client) {
	payload, err := json.Marshal(map[string]string{
		"connId": client.ConnID,
	})
	if err != nil {
		return
	}

	event := Event{
		Type:      EventTypeHello,
		Payload:   payload,
		UserID:    clientEventUserID(client),
		Timestamp: time.Now().UTC(),
	}
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
		s.unregister <- client
	}
}

func (s *Service) normalizeIncomingEvent(client *Client, event *Event) {
	if event == nil || client == nil {
		return
	}

	event.Type = strings.TrimSpace(event.Type)
	event.RoomID = strings.TrimSpace(event.RoomID)
	event.UserID = clientEventUserID(client)
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	} else {
		event.Timestamp = event.Timestamp.UTC()
	}
}

func (s *Service) normalizeOutgoingEvent(event *Event) {
	if event == nil {
		return
	}

	event.Type = strings.TrimSpace(event.Type)
	event.RoomID = strings.TrimSpace(event.RoomID)
	event.UserID = strings.TrimSpace(event.UserID)
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	} else {
		event.Timestamp = event.Timestamp.UTC()
	}
}

func resolveIdentityID(identity ConnectIdentity) (string, bool) {
	switch identity.Type {
	case IdentityTypeUser:
		userID := strings.TrimSpace(identity.UserID)
		if userID == "" {
			return "", false
		}
		return "user:" + userID, true
	case IdentityTypeGuest:
		participantID := strings.TrimSpace(identity.ParticipantID)
		if participantID == "" {
			return "", false
		}
		return "guest:" + participantID, true
	default:
		return "", false
	}
}

func clientEventUserID(client *Client) string {
	if client == nil {
		return ""
	}
	if strings.TrimSpace(client.UserID) != "" {
		return strings.TrimSpace(client.UserID)
	}
	return strings.TrimSpace(client.ParticipantID)
}

func (s *Service) removeClientFromRoomLocked(client *Client) (string, bool) {
	if client == nil || strings.TrimSpace(client.RoomID) == "" {
		return "", false
	}

	roomID := strings.TrimSpace(client.RoomID)
	if members, ok := s.roomClients[roomID]; ok {
		delete(members, client)
		if len(members) == 0 {
			delete(s.roomClients, roomID)
		}
	}

	presenceLeft := false
	if presence, ok := s.roomPresence[roomID]; ok {
		count := presence[client.IdentityID]
		if count <= 1 {
			delete(presence, client.IdentityID)
			presenceLeft = true
		} else {
			presence[client.IdentityID] = count - 1
		}

		if len(presence) == 0 {
			delete(s.roomPresence, roomID)
		}
	}

	client.RoomID = ""
	return roomID, presenceLeft
}
