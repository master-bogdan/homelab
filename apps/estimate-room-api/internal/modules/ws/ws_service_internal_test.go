package ws

import (
	"encoding/json"
	"testing"
	"time"
)

func TestServiceSendToUser_DeliversToMatchingIdentity(t *testing.T) {
	service := NewService(nil, "test")

	client := &Client{
		ConnID:       "conn-1",
		IdentityType: IdentityTypeUser,
		IdentityID:   "user:user-123",
		UserID:       "user-123",
		Send:         make(chan []byte, 1),
	}

	service.register <- client
	waitForRegisteredClient(t, service, client.ConnID)

	if err := service.SendToUser("user-123", Event{Type: "REWARD_READY"}); err != nil {
		t.Fatalf("expected send to user to succeed: %v", err)
	}

	select {
	case raw := <-client.Send:
		event := Event{}
		if err := json.Unmarshal(raw, &event); err != nil {
			t.Fatalf("failed to decode event: %v", err)
		}
		if event.Type != "REWARD_READY" {
			t.Fatalf("expected event type REWARD_READY, got %s", event.Type)
		}
		if event.UserID != "user-123" {
			t.Fatalf("expected event user id user-123, got %s", event.UserID)
		}
	case <-time.After(time.Second):
		t.Fatal("expected event to be delivered to client")
	}

	service.unregister <- client
}

func TestServiceSendToIdentity_IgnoresMissingIdentity(t *testing.T) {
	service := NewService(nil, "test")

	if err := service.SendToIdentity("user:missing", Event{Type: "NO_CLIENTS"}); err != nil {
		t.Fatalf("expected missing identity send to be a no-op, got %v", err)
	}
}

func waitForRegisteredClient(t *testing.T, service *Service, connID string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		service.mu.RLock()
		_, ok := service.connClients[connID]
		service.mu.RUnlock()
		if ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for client %s registration", connID)
}
