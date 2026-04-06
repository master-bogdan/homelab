package e2e

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/master-bogdan/estimate-room-api/internal/app"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

const e2eWSOrigin = "http://frontend.test"

type e2eApp struct {
	server      *httptest.Server
	db          *bun.DB
	clientID    string
	redirectURI string
}

func setupE2EApp(t *testing.T) *e2eApp {
	t.Helper()

	db := testutils.SetupTestDB(t)
	resetE2EDB(t, db)

	router := chi.NewRouter()
	cfg := &config.Config{}
	cfg.Frontend.BaseURL = "http://localhost:5173"
	cfg.Server.PasetoSymmetricKey = testutils.TestTokenKey
	cfg.Server.Issuer = testutils.TestIssuer
	cfg.Server.WebSocketAllowedOrigins = e2eWSOrigin

	backgroundCtx, cancel := context.WithCancel(context.Background())
	application := app.AppDeps{
		DB:                 db,
		Redis:              nil,
		Cfg:                cfg,
		Router:             router,
		IsGracefulShutdown: &atomic.Bool{},
		WsServer:           nil,
	}
	if err := application.SetupApp(backgroundCtx); err != nil {
		t.Fatalf("failed to set up app: %v", err)
	}

	server := httptest.NewServer(router)
	redirectURI := "http://localhost:4081/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})

	t.Cleanup(func() {
		server.Close()
		cancel()
		_ = db.Close()
	})

	return &e2eApp{
		server:      server,
		db:          db,
		clientID:    clientID,
		redirectURI: redirectURI,
	}
}

func resetE2EDB(t *testing.T, db *bun.DB) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
			invitations,
			votes,
			task_rounds,
			tasks,
			team_members,
			teams,
			room_participants,
			rooms,
			auth_password_reset_tokens,
			oauth2_access_tokens,
			oauth2_refresh_tokens,
			oauth2_auth_codes,
			oauth2_oidc_sessions,
			users,
			oauth2_clients
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate test tables: %v", err)
	}
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (a *e2eApp) loginAndGetAccessToken(t *testing.T, email, password string) string {
	t.Helper()

	state := uuid.NewString()
	nonce := uuid.NewString()
	codeVerifier := "verifier-" + uuid.NewString()
	codeChallenge := generateCodeChallenge(codeVerifier)

	authorizePath := "/api/v1/oauth2/authorize?client_id=" + url.QueryEscape(a.clientID) +
		"&redirect_uri=" + url.QueryEscape(a.redirectURI) +
		"&response_type=code" +
		"&scopes=user" +
		"&state=" + url.QueryEscape(state) +
		"&code_challenge=" + url.QueryEscape(codeChallenge) +
		"&code_challenge_method=S256" +
		"&nonce=" + url.QueryEscape(nonce)
	authorizeURL := a.server.URL + authorizePath

	loginBody := map[string]string{
		"email":    email,
		"password": password,
		"continue": authorizePath,
	}
	bodyBytes, err := json.Marshal(loginBody)
	if err != nil {
		t.Fatalf("failed to encode login request: %v", err)
	}

	loginReq, err := http.NewRequest(http.MethodPost, a.server.URL+"/api/v1/auth/login", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("failed to build login request: %v", err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginClient := &http.Client{
		Transport: a.server.Client().Transport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	loginResp, err := loginClient.Do(loginReq)
	if err != nil {
		t.Fatalf("failed to call login endpoint: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(loginResp.Body)
		t.Fatalf("expected 200 from auth login, got %d: %s", loginResp.StatusCode, string(body))
	}

	authorizeReq, err := http.NewRequest(http.MethodGet, authorizeURL, nil)
	if err != nil {
		t.Fatalf("failed to build authorize request: %v", err)
	}
	for _, cookie := range loginResp.Cookies() {
		authorizeReq.AddCookie(cookie)
	}

	authorizeClient := &http.Client{
		Transport: a.server.Client().Transport,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	authorizeResp, err := authorizeClient.Do(authorizeReq)
	if err != nil {
		t.Fatalf("failed to call authorize endpoint: %v", err)
	}
	defer authorizeResp.Body.Close()

	if authorizeResp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(authorizeResp.Body)
		t.Fatalf("expected 302 from authorize, got %d: %s", authorizeResp.StatusCode, string(body))
	}

	location := authorizeResp.Header.Get("Location")
	redirectURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse authorize redirect: %v", err)
	}

	code := strings.TrimSpace(redirectURL.Query().Get("code"))
	if code == "" {
		t.Fatalf("expected auth code in authorize redirect URL, got %q", location)
	}

	tokenForm := url.Values{}
	tokenForm.Set("grant_type", "authorization_code")
	tokenForm.Set("code", code)
	tokenForm.Set("redirect_uri", a.redirectURI)
	tokenForm.Set("client_id", a.clientID)
	tokenForm.Set("code_verifier", codeVerifier)

	tokenReq, err := http.NewRequest(http.MethodPost, a.server.URL+"/api/v1/oauth2/token", strings.NewReader(tokenForm.Encode()))
	if err != nil {
		t.Fatalf("failed to build token request: %v", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := a.server.Client().Do(tokenReq)
	if err != nil {
		t.Fatalf("failed to call token endpoint: %v", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(tokenResp.Body)
		t.Fatalf("expected 200 from token endpoint, got %d: %s", tokenResp.StatusCode, string(body))
	}

	payload := oauth2dto.TokenResponseDTO{}
	if err := json.NewDecoder(tokenResp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode token response: %v", err)
	}
	if strings.TrimSpace(payload.AccessToken) == "" {
		t.Fatal("expected access token in token response")
	}

	return payload.AccessToken
}

func doJSONRequest(
	t *testing.T,
	client *http.Client,
	method, rawURL string,
	body string,
	accessToken string,
) *http.Response {
	t.Helper()

	var requestBody io.Reader
	if body != "" {
		requestBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, rawURL, requestBody)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if strings.TrimSpace(accessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	return resp
}

func decodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	defer resp.Body.Close()

	var payload T
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	return payload
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	return string(body)
}

func connectWS(t *testing.T, serverURL, accessToken string) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(serverURL, "http") + "/api/v1/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization": []string{"Bearer " + accessToken},
			"Origin":        []string{e2eWSOrigin},
		},
	})
	if err != nil {
		t.Fatalf("failed to connect websocket: %v", err)
	}

	return conn
}

func readUntilEvent(t *testing.T, conn *websocket.Conn, eventType string) ws.Event {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("failed to read websocket event %s: %v", eventType, err)
		}

		event := ws.Event{}
		if err := json.Unmarshal(data, &event); err != nil {
			t.Fatalf("failed to decode websocket event: %v", err)
		}

		if event.Type == eventType {
			return event
		}
	}
}

func writeEvent(t *testing.T, conn *websocket.Conn, event ws.Event) {
	t.Helper()

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal websocket event: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("failed to write websocket event: %v", err)
	}
}

func mustMarshalJSON(t *testing.T, value any) []byte {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("failed to marshal JSON payload: %v", err)
	}

	return data
}
