package tests

// import (
// 	"bytes"
// 	"crypto/sha256"
// 	"database/sql"
// 	"encoding/base64"
// 	"io"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"strings"
// 	"testing"
//
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/master-bogdan/clear-cash-api/config"
// 	"github.com/master-bogdan/clear-cash-api/internal/infra/db/postgresql/repositories"
// 	"github.com/master-bogdan/clear-cash-api/internal/modules/oauth2"
// 	test_utils "github.com/master-bogdan/clear-cash-api/internal/pkg/test"
// )
//
// var userID string
// var clientID string
// var sessionID string
//
// func setupTest(t *testing.T) (*fiber.App, *sql.DB, *config.Config) {
// 	app, db, cfg, err := test_utils.SetupTestApp()
// 	if err != nil {
// 		t.Fatalf("failed to setup test app: %v", err)
// 	}
//
// 	v1 := app.Group("/api/v1")
//
// 	// Create repos from db
// 	clientRepo := repositories.NewOauth2ClientRepository(db)
// 	authCodeRepo := repositories.NewOauth2AuthCodeRepository(db)
// 	userRepo := repositories.NewOauth2UserRepository(db)
// 	oidcSessionRepo := repositories.NewOauth2OidcSessionRepository(db)
// 	refreshTokenRepo := repositories.NewOauth2RefreshTokenRepository(db)
// 	accessTokenRepo := repositories.NewOauth2AccessTokenRepository(db)
//
// 	// Inject repos and tokenKey
// 	oauth2.NewOauth2Module(
// 		v1,
// 		cfg.Token.TokenKey,
// 		clientRepo,
// 		authCodeRepo,
// 		userRepo,
// 		oidcSessionRepo,
// 		refreshTokenRepo,
// 		accessTokenRepo,
// 		slog.Default(),
// 	)
//
// 	userID, clientID, sessionID, _, _, _, err = test_utils.SeedAuth(db)
// 	if err != nil {
// 		t.Fatalf("failed to seed auth: %v", err)
// 	}
//
// 	return app, db, cfg
// }
//
// func TestAuthorize_WithValidSession_ReturnsCode_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
// 		"&redirect_uri=http://localhost:4081"+
// 		"&response_type=code"+
// 		"&state=xyz123"+
// 		"&scopes=user"+
// 		"&nonce=nonce123"+
// 		"&code_challenge=abc123"+
// 		"&code_challenge_method=S256", nil)
// 	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
//
// 	resp, err := app.Test(req)
// 	if err != nil {
// 		t.Fatalf("failed to send request: %v", err)
// 	}
//
// 	if resp.StatusCode != fiber.StatusFound {
// 		t.Errorf("expected 302 Found, got %d", resp.StatusCode)
// 	}
//
// 	loc := resp.Header.Get("Location")
// 	if !strings.Contains(loc, "code=") {
// 		t.Errorf("expected Location header to contain auth code, got: %s", loc)
// 	}
// }
//
// func TestAuthorize_WithoutSession_RedirectsToLogin_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
// 		"&redirect_uri=http://localhost:4081"+
// 		"&response_type=code"+
// 		"&state=xyz123"+
// 		"&scopes=user"+
// 		"&nonce=nonce123"+
// 		"&code_challenge=abc123"+
// 		"&code_challenge_method=S256", nil)
//
// 	resp, err := app.Test(req)
// 	if err != nil {
// 		t.Fatalf("failed to send request: %v", err)
// 	}
//
// 	if resp.StatusCode != fiber.StatusFound {
// 		t.Errorf("expected 302 Found, got %d", resp.StatusCode)
// 	}
//
// 	loc := resp.Header.Get("Location")
// 	if !strings.HasPrefix(loc, "/api/v1/oauth2/login") {
// 		t.Errorf("expected redirect to login, got: %s", loc)
// 	}
// }
//
// func TestLogin_ValidUser_RedirectsWithCode_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	form := url.Values{}
// 	form.Set("email", "testuser@example.com")
// 	form.Set("password", "fakehash")
// 	form.Set("client_id", clientID)
// 	form.Set("redirect_uri", "http://localhost:4081")
// 	form.Set("response_type", "code")
// 	form.Set("scopes", "user")
// 	form.Set("state", "xyz123")
// 	form.Set("code_challenge", "abc123")
// 	form.Set("code_challenge_method", "S256")
// 	form.Set("nonce", "nonce123")
//
// 	req := httptest.NewRequest("POST", "/api/v1/oauth2/login", strings.NewReader(form.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	resp, err := app.Test(req)
// 	if err != nil {
// 		t.Fatalf("failed to send request: %v", err)
// 	}
//
// 	if resp.StatusCode != fiber.StatusFound {
// 		t.Errorf("expected 302 redirect, got %d", resp.StatusCode)
// 	}
//
// 	loc := resp.Header.Get("Location")
// 	if !strings.Contains(loc, "code=") {
// 		t.Errorf("expected redirect to contain auth code, got: %s", loc)
// 	}
// }
//
// func TestShowLoginForm_ReturnsHtml_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	req := httptest.NewRequest("GET", "/api/v1/oauth2/login?client_id="+clientID+
// 		"&redirect_uri=http://localhost:4081"+
// 		"&response_type=code"+
// 		"&state=xyz123"+
// 		"&scopes=user"+
// 		"&nonce=nonce123"+
// 		"&code_challenge=abc123"+
// 		"&code_challenge_method=S256", nil)
//
// 	resp, err := app.Test(req)
// 	if err != nil {
// 		t.Fatalf("failed to send request: %v", err)
// 	}
//
// 	if resp.StatusCode != fiber.StatusOK {
// 		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
// 	}
//
// 	ct := resp.Header.Get("Content-Type")
// 	if !strings.HasPrefix(ct, "text/html") {
// 		t.Errorf("expected text/html response, got: %s", ct)
// 	}
// }
//
// func generateCodeChallenge(verifier string) string {
// 	hash := sha256.Sum256([]byte(verifier))
// 	return base64.RawURLEncoding.EncodeToString(hash[:])
// }
//
// func TestToken_AuthorizationCodeFlow_ReturnsTokens_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	// Step 0: generate PKCE values
// 	codeVerifier := "testverifier1234567890"
// 	codeChallenge := generateCodeChallenge(codeVerifier)
// 	state := "xyz123"
//
// 	// Step 1: login / authorize
// 	formLogin := url.Values{}
// 	formLogin.Set("email", "testuser@example.com")
// 	formLogin.Set("password", "fakehash")
// 	formLogin.Set("client_id", clientID)
// 	formLogin.Set("redirect_uri", "http://localhost:4081")
// 	formLogin.Set("response_type", "code")
// 	formLogin.Set("scopes", "user")
// 	formLogin.Set("state", state)
// 	formLogin.Set("code_challenge", codeChallenge)
// 	formLogin.Set("code_challenge_method", "S256")
// 	formLogin.Set("nonce", "nonce123")
//
// 	reqLogin := httptest.NewRequest("POST", "/api/v1/oauth2/login", strings.NewReader(formLogin.Encode()))
// 	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	reqLogin.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})
//
// 	respLogin, err := app.Test(reqLogin)
// 	if err != nil {
// 		t.Fatalf("failed login request: %v", err)
// 	}
//
// 	if respLogin.StatusCode != fiber.StatusFound {
// 		t.Fatalf("expected 302 Found, got %d", respLogin.StatusCode)
// 	}
//
// 	loc := respLogin.Header.Get("Location")
// 	u, _ := url.Parse(loc)
// 	code := u.Query().Get("code")
// 	if code == "" {
// 		t.Fatalf("expected auth code in redirect URL")
// 	}
//
// 	// Step 2: exchange code for tokens
// 	formToken := url.Values{}
// 	formToken.Set("grant_type", "authorization_code")
// 	formToken.Set("code", code)
// 	formToken.Set("redirect_uri", "http://localhost:4081")
// 	formToken.Set("client_id", clientID)
// 	formToken.Set("code_verifier", codeVerifier) // must match PKCE verifier
//
// 	reqToken := httptest.NewRequest("POST", "/api/v1/oauth2/token", strings.NewReader(formToken.Encode()))
// 	reqToken.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	respToken, err := app.Test(reqToken)
// 	if err != nil {
// 		t.Fatalf("failed token request: %v", err)
// 	}
//
// 	if respToken.StatusCode != fiber.StatusOK {
// 		t.Fatalf("expected 200 OK, got %d", respToken.StatusCode)
// 	}
//
// 	respBody, err := io.ReadAll(respToken.Body)
// 	if err != nil {
// 		t.Fatalf("failed to read response body: %v", err)
// 	}
// 	respToken.Body.Close()
//
// 	bodyStr := string(respBody)
//
// 	if !strings.Contains(bodyStr, "access_token") {
// 		t.Errorf("expected response to contain access_token, got %s", bodyStr)
// 	}
// 	if !strings.Contains(bodyStr, "refresh_token") {
// 		t.Errorf("expected response to contain refresh_token, got %s", bodyStr)
// 	}
// 	if !strings.Contains(bodyStr, "id_token") {
// 		t.Errorf("expected response to contain id_token, got %s", bodyStr)
// 	}
// }
//
// func TestAuthorize_InvalidQueryParams_ReturnsBadRequest_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?invalid=1", nil)
// 	resp, _ := app.Test(req)
// 	if resp.StatusCode != fiber.StatusBadRequest {
// 		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
// 	}
// }
//
// func TestLogin_InvalidJsonBody_ReturnsBadRequest_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	req := httptest.NewRequest("POST", "/api/v1/oauth2/login", bytes.NewBuffer([]byte("notjson")))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, _ := app.Test(req)
// 	if resp.StatusCode != fiber.StatusBadRequest {
// 		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
// 	}
// }
//
// func TestGetTokens_UnsupportedGrantType_ReturnsBadRequest_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	form := url.Values{}
// 	form.Set("grant_type", "foobar")
//
// 	req := httptest.NewRequest("POST", "/api/v1/oauth2/token", strings.NewReader(form.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	resp, _ := app.Test(req)
// 	if resp.StatusCode != fiber.StatusBadRequest {
// 		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
// 	}
// }
//
// func TestAuthorize_InvalidRedirectURI_ReturnsBadRequest_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
// 		"&redirect_uri=http://malicious.com"+
// 		"&response_type=code"+
// 		"&code_challenge=abc123"+
// 		"&code_challenge_method=S256"+
// 		"&scope=user", nil)
//
// 	resp, _ := app.Test(req)
// 	if resp.StatusCode != fiber.StatusBadRequest {
// 		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
// 	}
// }
//
// func TestLogin_MissingRequiredField_ReturnsBadRequest_E2E(t *testing.T) {
// 	app, db, _ := setupTest(t)
// 	defer db.Close()
//
// 	form := url.Values{}
// 	form.Set("email", "testuser@example.com")
// 	form.Set("password", "fakehash")
// 	form.Set("client_id", clientID)
// 	form.Set("redirect_uri", "http://localhost:4081")
// 	form.Set("response_type", "code")
// 	form.Set("scopes", "user")
// 	form.Set("state", "xyz123")
// 	form.Set("nonce", "nonce123")
//
// 	req := httptest.NewRequest("POST", "/api/v1/oauth2/login", strings.NewReader(form.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
// 	resp, _ := app.Test(req)
// 	if resp.StatusCode != fiber.StatusBadRequest {
// 		t.Errorf("expected 400 Bad Request, got %d", resp.StatusCode)
// 	}
// }
