package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

const (
	pingInterval = 30 * time.Second
	pingTimeout  = 10 * time.Second
)

type Client struct {
	Conn      *websocket.Conn
	ClientID  string
	ChannelID string
	Send      chan []byte
	Gateway   Gateway
}

type ClientInfo struct {
	ClientID  string
	ChannelID string
}

type Gateway interface {
	HandleConnection(w http.ResponseWriter, r *http.Request)
	OnConnect(client ClientInfo)
	OnDisconnect(client ClientInfo)
	OnMessage(client ClientInfo, message []byte)
	OnError(client ClientInfo, err error)
}

type Manager struct {
	clients        map[*Client]bool
	channelClients map[string]map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	mu             sync.RWMutex
	server         *WsServer
	channel        string
}

func NewManager(server *WsServer, channel string) *Manager {
	m := &Manager{
		clients:        make(map[*Client]bool),
		channelClients: make(map[string]map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		server:         server,
		channel:        channel,
	}

	go m.run()

	server.Subscribe(channel, func(data []byte) {
		var msg map[string]any
		if err := json.Unmarshal(data, &msg); err == nil {
			if channelID, ok := msg["channelID"].(string); ok {
				m.BroadcastToChannel(channelID, data)
			}
		}
	})

	return m
}

func clientInfo(client *Client) ClientInfo {
	return ClientInfo{
		ClientID:  client.ClientID,
		ChannelID: client.ChannelID,
	}
}

func (m *Manager) run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true

			if m.channelClients[client.ChannelID] == nil {
				m.channelClients[client.ChannelID] = make(map[*Client]bool)
			}
			m.channelClients[client.ChannelID][client] = true
			m.mu.Unlock()

			if client.Gateway != nil {
				client.Gateway.OnConnect(clientInfo(client))
			}

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				delete(m.channelClients[client.ChannelID], client)
				if len(m.channelClients[client.ChannelID]) == 0 {
					delete(m.channelClients, client.ChannelID)
				}
				close(client.Send)
			}
			m.mu.Unlock()

			if client.Gateway != nil {
				client.Gateway.OnDisconnect(clientInfo(client))
			}
		}
	}
}

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request, channelID, clientID string, gateway Gateway) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("WebSocket accept error: %v", err)
		return
	}

	client := &Client{
		Conn:      conn,
		ClientID:  clientID,
		ChannelID: channelID,
		Send:      make(chan []byte, 256),
		Gateway:   gateway,
	}

	m.register <- client

	go m.writeHandler(client)
	m.readHandler(client)
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
			if client.Gateway != nil {
				client.Gateway.OnError(clientInfo(client), err)
			}
			break
		}
		if client.Gateway != nil {
			client.Gateway.OnMessage(clientInfo(client), message)
		}
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
				if client.Gateway != nil {
					client.Gateway.OnError(clientInfo(client), err)
				}
				return
			}
		case <-pingTicker.C:
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
			err := client.Conn.Ping(ctx)
			cancel()
			if err != nil {
				if client.Gateway != nil {
					client.Gateway.OnError(clientInfo(client), err)
				}
				client.Conn.Close(websocket.StatusNormalClosure, "ping timeout")
				return
			}
		}
	}
}

func (m *Manager) BroadcastToChannel(channelID string, data []byte) {
	m.mu.RLock()
	clients := make([]*Client, 0, len(m.channelClients[channelID]))
	for client := range m.channelClients[channelID] {
		clients = append(clients, client)
	}
	m.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			m.unregister <- client
		}
	}
}

func (m *Manager) Broadcast(message any) error {
	return m.server.Publish(m.channel, message)
}

func (m *Manager) ClientIDs(channelID string) []string {
	m.mu.RLock()
	channelClients := m.channelClients[channelID]
	uniqueIDs := make(map[string]struct{}, len(channelClients))
	for client := range channelClients {
		if client.ClientID != "" {
			uniqueIDs[client.ClientID] = struct{}{}
		}
	}
	m.mu.RUnlock()

	ids := make([]string, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		ids = append(ids, id)
	}
	return ids
}
