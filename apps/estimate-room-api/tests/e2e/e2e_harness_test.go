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

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/master-bogdan/estimate-room-api/internal/app"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

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
	cfg.Server.PasetoSymmetricKey = testutils.TestTokenKey
	cfg.Server.Issuer = testutils.TestIssuer

	backgroundCtx, cancel := context.WithCancel(context.Background())
	application := app.AppDeps{
		DB:                 db,
		Redis:              nil,
		Cfg:                cfg,
		Router:             router,
		IsGracefulShutdown: &atomic.Bool{},
		WsServer:           nil,
	}
	application.SetupApp(backgroundCtx)

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

	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	form.Set("client_id", a.clientID)
	form.Set("redirect_uri", a.redirectURI)
	form.Set("response_type", "code")
	form.Set("scopes", "user")
	form.Set("state", state)
	form.Set("code_challenge", codeChallenge)
	form.Set("code_challenge_method", "S256")
	form.Set("nonce", nonce)

	loginReq, err := http.NewRequest(http.MethodPost, a.server.URL+"/api/v1/oauth2/login", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("failed to build login request: %v", err)
	}
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	if loginResp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(loginResp.Body)
		t.Fatalf("expected 302 from login, got %d: %s", loginResp.StatusCode, string(body))
	}

	location := loginResp.Header.Get("Location")
	redirectURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse login redirect: %v", err)
	}

	code := strings.TrimSpace(redirectURL.Query().Get("code"))
	if code == "" {
		t.Fatalf("expected auth code in redirect URL, got %q", location)
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
