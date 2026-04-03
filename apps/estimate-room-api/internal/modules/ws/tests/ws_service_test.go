package tests

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
)

func TestServiceConnect_AllowsConfiguredOrigin(t *testing.T) {
	service := ws.NewService(nil, "test")
	service.SetOriginPatterns([]string{"http://frontend.test"})

	req := newWebSocketRequest("http://api.test/ws", "http://frontend.test")
	writer := newHijackableResponseWriter(t)
	defer writer.Close()

	done := make(chan struct{})
	go func() {
		service.Connect(writer, req, ws.ConnectIdentity{
			Type:   ws.IdentityTypeUser,
			UserID: "user-123",
		})
		close(done)
	}()

	writer.CloseClient()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("expected websocket connect to finish after client close")
	}

	if writer.statusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected 101 Switching Protocols for allowed origin, got %d", writer.statusCode)
	}
}

func TestServiceConnect_RejectsUnauthorizedOrigin(t *testing.T) {
	service := ws.NewService(nil, "test")
	service.SetOriginPatterns([]string{"http://frontend.test"})

	req := newWebSocketRequest("http://api.test/ws", "http://evil.test")
	rr := httptest.NewRecorder()

	service.Connect(rr, req, ws.ConnectIdentity{
		Type:   ws.IdentityTypeUser,
		UserID: "user-123",
	})

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden for unauthorized origin, got %d", rr.Code)
	}
}

func newWebSocketRequest(targetURL, origin string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, targetURL, nil)
	req.Host = strings.TrimPrefix(strings.TrimPrefix(targetURL, "http://"), "https://")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	req.Header.Set("Origin", origin)
	return req
}

type hijackableResponseWriter struct {
	header     http.Header
	statusCode int
	serverConn net.Conn
	clientConn net.Conn
}

func newHijackableResponseWriter(t *testing.T) *hijackableResponseWriter {
	t.Helper()

	serverConn, clientConn := net.Pipe()
	return &hijackableResponseWriter{
		header:     make(http.Header),
		serverConn: serverConn,
		clientConn: clientConn,
	}
}

func (w *hijackableResponseWriter) Header() http.Header {
	return w.header
}

func (w *hijackableResponseWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *hijackableResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *hijackableResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.serverConn, bufio.NewReadWriter(bufio.NewReader(w.serverConn), bufio.NewWriter(w.serverConn)), nil
}

func (w *hijackableResponseWriter) CloseClient() {
	_ = w.clientConn.Close()
}

func (w *hijackableResponseWriter) Close() {
	_ = w.clientConn.Close()
	_ = w.serverConn.Close()
}
