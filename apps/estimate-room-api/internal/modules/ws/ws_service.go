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
	Conn   *websocket.Conn
	UserID string
	ConnID string
	Send   chan []byte
}

type Service struct {
	clients       map[*Client]bool
	userClients   map[string]map[*Client]bool
	register      chan *Client
	unregister    chan *Client
	mu            sync.RWMutex
	server        PubSub
	channel       string
	subscriptions map[string][]EventHandler
}

func NewService(server PubSub, channel string) *Service {
	s := &Service{
		clients:       make(map[*Client]bool),
		userClients:   make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		server:        server,
		channel:       channel,
		subscriptions: make(map[string][]EventHandler),
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
	return ClientInfo{
		UserID: client.UserID,
		ConnID: client.ConnID,
	}
}

func (s *Service) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true

			if s.userClients[client.UserID] == nil {
				s.userClients[client.UserID] = make(map[*Client]bool)
			}
			s.userClients[client.UserID][client] = true
			s.mu.Unlock()

			s.logConnect(clientInfo(client))

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				delete(s.userClients[client.UserID], client)
				if len(s.userClients[client.UserID]) == 0 {
					delete(s.userClients, client.UserID)
				}
				close(client.Send)
			}
			s.mu.Unlock()

			s.logDisconnect(clientInfo(client))
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

func (s *Service) Connect(w http.ResponseWriter, r *http.Request, userID string) {
	if strings.TrimSpace(userID) == "" {
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
		Conn:   conn,
		UserID: userID,
		ConnID: connID,
		Send:   make(chan []byte, 256),
	}

	s.disconnectUserConnections(userID)

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

func (s *Service) Broadcast(message any) error {
	if s.server == nil {
		return errors.New("ws server is not initialized")
	}
	return s.server.Publish(s.channel, message)
}

func (s *Service) disconnectUserConnections(userID string) {
	s.mu.RLock()
	existing := make([]*Client, 0, len(s.userClients[userID]))
	for client := range s.userClients[userID] {
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

func (s *Service) logConnect(info ClientInfo) {
	logger.L().Info("ws connected", "user_id", info.UserID, "conn_id", info.ConnID)
}

func (s *Service) logDisconnect(info ClientInfo) {
	logger.L().Info("ws disconnected", "user_id", info.UserID, "conn_id", info.ConnID)
}

func (s *Service) logError(info ClientInfo, err error) {
	logger.L().Error("ws error", "user_id", info.UserID, "conn_id", info.ConnID, "err", err)
}

func (s *Service) sendHello(client *Client) {
	payload, err := json.Marshal(map[string]string{
		"connId": client.ConnID,
		"userId": client.UserID,
	})
	if err != nil {
		return
	}

	event := Event{
		Type:    "HELLO",
		Payload: payload,
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
