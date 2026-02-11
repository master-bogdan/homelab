package ws

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
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
