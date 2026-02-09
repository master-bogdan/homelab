package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
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

type ClientInfo struct {
	UserID string
	ConnID string
}

type EventHandler func(ClientInfo, Event)

type Manager struct {
	clients       map[*Client]bool
	userClients   map[string]map[*Client]bool
	register      chan *Client
	unregister    chan *Client
	mu            sync.RWMutex
	server        *WsServer
	channel       string
	subscriptions map[string][]EventHandler
}

func NewManager(server *WsServer, channel string) *Manager {
	m := &Manager{
		clients:       make(map[*Client]bool),
		userClients:   make(map[string]map[*Client]bool),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		server:        server,
		channel:       channel,
		subscriptions: make(map[string][]EventHandler),
	}

	go m.run()

	server.Subscribe(channel, func(data []byte) {
		m.broadcastRaw(data)
	})

	return m
}

func clientInfo(client *Client) ClientInfo {
	return ClientInfo{
		UserID: client.UserID,
		ConnID: client.ConnID,
	}
}

func (m *Manager) run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true

			if m.userClients[client.UserID] == nil {
				m.userClients[client.UserID] = make(map[*Client]bool)
			}
			m.userClients[client.UserID][client] = true
			m.mu.Unlock()

			m.logConnect(clientInfo(client))

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				delete(m.userClients[client.UserID], client)
				if len(m.userClients[client.UserID]) == 0 {
					delete(m.userClients, client.UserID)
				}
				close(client.Send)
			}
			m.mu.Unlock()

			m.logDisconnect(clientInfo(client))
		}
	}
}

func (m *Manager) Subscribe(eventType string, handler EventHandler) {
	if strings.TrimSpace(eventType) == "" || handler == nil {
		return
	}
	m.mu.Lock()
	m.subscriptions[eventType] = append(m.subscriptions[eventType], handler)
	m.mu.Unlock()
}

func (m *Manager) Connect(w http.ResponseWriter, r *http.Request, userID string) {
	if strings.TrimSpace(userID) == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if !isOriginAllowed(r) {
		http.Error(w, "invalid origin", http.StatusForbidden)
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

	m.kickUserConnections(userID)

	m.register <- client

	go m.writeHandler(client)
	m.sendHello(client)
	m.readHandler(client)
}

func isOriginAllowed(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}

	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	originHost := hostOnly(originURL.Host)
	requestHost := hostOnly(r.Host)
	if originHost == "" || requestHost == "" {
		return false
	}

	return strings.EqualFold(originHost, requestHost)
}

func hostOnly(hostport string) string {
	if hostport == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(hostport); err == nil {
		return host
	}
	return hostport
}

func (m *Manager) readHandler(client *Client) {
	defer func() {
		m.unregister <- client
		client.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	ctx := context.Background()
	for {
		_, message, err := client.Conn.Read(ctx)
		if err != nil {
			m.logError(clientInfo(client), err)
			break
		}
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			m.logError(clientInfo(client), err)
			continue
		}
		if strings.TrimSpace(event.Type) == "" {
			m.logError(clientInfo(client), errors.New("missing event type"))
			continue
		}
		m.dispatchEvent(clientInfo(client), event)
	}
}

func (m *Manager) writeHandler(client *Client) {
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
				m.logError(clientInfo(client), err)
				return
			}
		case <-pingTicker.C:
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
			err := client.Conn.Ping(ctx)
			cancel()
			if err != nil {
				m.logError(clientInfo(client), err)
				client.Conn.Close(websocket.StatusNormalClosure, "ping timeout")
				return
			}
		}
	}
}

func (m *Manager) broadcastRaw(data []byte) {
	m.mu.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for client := range m.clients {
		clients = append(clients, client)
	}
	m.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			m.unregister <- client
		}
	}
}

func (m *Manager) Broadcast(message any) error {
	return m.server.Publish(m.channel, message)
}

func (m *Manager) kickUserConnections(userID string) {
	m.mu.RLock()
	existing := make([]*Client, 0, len(m.userClients[userID]))
	for client := range m.userClients[userID] {
		existing = append(existing, client)
	}
	m.mu.RUnlock()

	for _, client := range existing {
		_ = client.Conn.Close(websocket.StatusPolicyViolation, "another connection opened")
	}
}

func (m *Manager) dispatchEvent(info ClientInfo, event Event) {
	m.mu.RLock()
	handlers := append([]EventHandler(nil), m.subscriptions[event.Type]...)
	m.mu.RUnlock()

	for _, handler := range handlers {
		handler(info, event)
	}
}

func (m *Manager) logConnect(info ClientInfo) {
	logger.L().Info("ws connected", "user_id", info.UserID, "conn_id", info.ConnID)
}

func (m *Manager) logDisconnect(info ClientInfo) {
	logger.L().Info("ws disconnected", "user_id", info.UserID, "conn_id", info.ConnID)
}

func (m *Manager) logError(info ClientInfo, err error) {
	logger.L().Error("ws error", "user_id", info.UserID, "conn_id", info.ConnID, "err", err)
}

func (m *Manager) sendHello(client *Client) {
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
		m.unregister <- client
	}
}
