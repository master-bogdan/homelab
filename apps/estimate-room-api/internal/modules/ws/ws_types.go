package ws

import (
	"encoding/json"
	"time"
)

const (
	EventTypeHello = "HELLO"
)

type IdentityType string

const (
	IdentityTypeUser  IdentityType = "USER"
	IdentityTypeGuest IdentityType = "GUEST"
)

type Event struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	RoomID    string          `json:"roomId,omitempty"`
	UserID    string          `json:"userId,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

type ConnectIdentity struct {
	Type          IdentityType
	UserID        string
	ParticipantID string
}

type ClientInfo struct {
	ConnID        string
	IdentityType  IdentityType
	IdentityID    string
	UserID        string
	ParticipantID string
}

type EventHandler func(ClientInfo, Event)

type PubSub interface {
	Subscribe(channel string, onMessage func([]byte))
	Publish(channel string, message any) error
}
