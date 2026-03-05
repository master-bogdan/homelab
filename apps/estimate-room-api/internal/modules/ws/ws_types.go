package ws

import (
	"encoding/json"
	"time"
)

const (
	EventTypeHello = "HELLO"
)

type Event struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	RoomID    string          `json:"roomId,omitempty"`
	UserID    string          `json:"userId,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type ClientInfo struct {
	UserID string
	ConnID string
}

type EventHandler func(ClientInfo, Event)

type PubSub interface {
	Subscribe(channel string, onMessage func([]byte))
	Publish(channel string, message any) error
}
