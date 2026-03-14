package ws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestServiceConnect_AllowsConfiguredOrigin(t *testing.T) {
	service := NewService(nil, "test")
	service.SetOriginPatterns([]string{"http://frontend.test"})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.Connect(w, r, ConnectIdentity{
			Type:   IdentityTypeUser,
			UserID: "user-123",
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(server.URL), &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Origin": []string{"http://frontend.test"},
		},
	})
	if err != nil {
		t.Fatalf("expected websocket dial to succeed for allowed origin: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")
}

func TestServiceConnect_RejectsUnauthorizedOrigin(t *testing.T) {
	service := NewService(nil, "test")
	service.SetOriginPatterns([]string{"http://frontend.test"})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service.Connect(w, r, ConnectIdentity{
			Type:   IdentityTypeUser,
			UserID: "user-123",
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, resp, err := websocket.Dial(ctx, wsURL(server.URL), &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Origin": []string{"http://evil.test"},
		},
	})
	if err == nil {
		t.Fatal("expected websocket dial to fail for unauthorized origin")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for unauthorized origin, got %v", resp)
	}
}

func wsURL(serverURL string) string {
	return "ws" + strings.TrimPrefix(serverURL, "http")
}
