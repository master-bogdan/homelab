package tests

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	authdto "github.com/master-bogdan/estimate-room-api/internal/modules/auth/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupAuthTest(t *testing.T) (*chi.Mux, *bun.DB, string, string) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	router := chi.NewRouter()
	sessionService := oauth2.NewAuthServiceFromDB(testutils.TestTokenKey, db)

	router.Route("/api/v1", func(r chi.Router) {
		userService := users.NewUsersService(usersrepositories.NewUserRepository(db))

		oauth2Module := oauth2.NewOauth2Module(oauth2.Oauth2ModuleDeps{
			Router:          r,
			DB:              db,
			TokenKey:        testutils.TestTokenKey,
			Issuer:          testutils.TestIssuer,
			UserService:     userService,
			AuthService:     sessionService,
			FrontendBaseURL: "http://localhost:5173",
		})

		auth.NewAuthModule(auth.AuthModuleDeps{
			Router:         r,
			DB:             db,
			UserService:    userService,
			Oauth2Service:  oauth2Module.Service,
			SessionService: oauth2Module.AuthService,
		})
	})

	redirectURI := "http://localhost:4081/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})

	return router, db, clientID, redirectURI
}

func continueURL(clientID, redirectURI, state, nonce, codeChallenge string) string {
	return "http://api.estimateroom.test/api/v1/oauth2/authorize?client_id=" + url.QueryEscape(clientID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code" +
		"&scopes=user" +
		"&state=" + url.QueryEscape(state) +
		"&code_challenge=" + url.QueryEscape(codeChallenge) +
		"&code_challenge_method=S256" +
		"&nonce=" + url.QueryEscape(nonce)
}

func TestLogin_ReturnsSessionCookieAndAuthenticatedPayload(t *testing.T) {
	router, db, clientID, redirectURI := setupAuthTest(t)
	defer db.Close()

	testutils.SeedUser(t, db, "testuser@example.com", "password123")

	body := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
		"continue": continueURL(clientID, redirectURI, "state-login", "nonce-login", "challenge-login"),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var response authdto.SessionResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Authenticated || response.User == nil {
		t.Fatalf("expected authenticated response, got %#v", response)
	}
	if response.User.Email == nil || *response.User.Email != "testuser@example.com" {
		t.Fatalf("expected user email in response, got %#v", response.User)
	}

	var sessionCookie *http.Cookie
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == oauth2.SessionCookieName {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil || sessionCookie.Value == "" {
		t.Fatalf("expected session cookie to be set")
	}
}

func TestRegister_DuplicateEmail_ReturnsConflict(t *testing.T) {
	router, db, clientID, redirectURI := setupAuthTest(t)
	defer db.Close()

	testutils.SeedUser(t, db, "taken@example.com", "password123")

	body := map[string]string{
		"email":       "taken@example.com",
		"password":    "password123",
		"displayName": "Taken User",
		"continue":    continueURL(clientID, redirectURI, "state-register", "nonce-register", "challenge-register"),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict, got %d", rr.Code)
	}
}

func TestRegister_PersistsProfessionalProfile(t *testing.T) {
	router, db, clientID, redirectURI := setupAuthTest(t)
	defer db.Close()

	body := map[string]string{
		"email":        "profile@example.com",
		"password":     "password123",
		"displayName":  "Profile User",
		"organization": "Acme Corp",
		"occupation":   "Developer",
		"continue":     continueURL(clientID, redirectURI, "state-profile", "nonce-profile", "challenge-profile"),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}

	var response authdto.SessionResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.User == nil {
		t.Fatalf("expected user payload in response")
	}
	if response.User.Organization == nil || *response.User.Organization != "Acme Corp" {
		t.Fatalf("expected organization in response, got %#v", response.User)
	}
	if response.User.Occupation == nil || *response.User.Occupation != "Developer" {
		t.Fatalf("expected occupation in response, got %#v", response.User)
	}

	userRepo := usersrepositories.NewUserRepository(db)
	user, err := userRepo.FindByEmail("profile@example.com")
	if err != nil {
		t.Fatalf("failed to load persisted user: %v", err)
	}
	if user.Organization == nil || *user.Organization != "Acme Corp" {
		t.Fatalf("expected organization to be persisted, got %#v", user.Organization)
	}
	if user.Occupation == nil || *user.Occupation != "Developer" {
		t.Fatalf("expected occupation to be persisted, got %#v", user.Occupation)
	}
}

func TestGetSession_WithSessionCookie_ReturnsAuthenticatedUser(t *testing.T) {
	router, db, clientID, redirectURI := setupAuthTest(t)
	defer db.Close()

	testutils.SeedUser(t, db, "session@example.com", "password123")

	loginBody := map[string]string{
		"email":    "session@example.com",
		"password": "password123",
		"continue": continueURL(clientID, redirectURI, "state-session", "nonce-session", "challenge-session"),
	}
	payload, err := json.Marshal(loginBody)
	if err != nil {
		t.Fatalf("failed to marshal login request: %v", err)
	}

	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(payload))
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	router.ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", loginRR.Code)
	}

	sessionReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	for _, cookie := range loginRR.Result().Cookies() {
		sessionReq.AddCookie(cookie)
	}

	sessionRR := httptest.NewRecorder()
	router.ServeHTTP(sessionRR, sessionReq)

	if sessionRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", sessionRR.Code)
	}

	var response authdto.SessionResponse
	if err := json.NewDecoder(sessionRR.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode session response: %v", err)
	}
	if !response.Authenticated || response.User == nil {
		t.Fatalf("expected authenticated session, got %#v", response)
	}
}

func TestValidateResetPasswordToken_ExpiredToken_ReturnsExpired(t *testing.T) {
	router, db, _, _ := setupAuthTest(t)
	defer db.Close()

	userID := testutils.SeedUser(t, db, "expired@example.com", "password123")
	rawToken := "expired-reset-token"
	seedResetToken(t, db, userID, rawToken, time.Now().Add(-time.Hour), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/reset-password/validate?token="+url.QueryEscape(rawToken), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var response authdto.ResetPasswordValidationResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode validation response: %v", err)
	}
	if response.Valid || response.Reason != "expired" {
		t.Fatalf("expected expired token response, got %#v", response)
	}
}

func TestResetPassword_RevokesExistingSession(t *testing.T) {
	router, db, clientID, _ := setupAuthTest(t)
	defer db.Close()

	userID := testutils.SeedUser(t, db, "reset@example.com", "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-reset-existing")
	rawToken := "valid-reset-token"
	seedResetToken(t, db, userID, rawToken, time.Now().Add(time.Hour), nil)

	body := map[string]string{
		"token":    rawToken,
		"password": "newpassword123",
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal reset request: %v", err)
	}

	resetReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(payload))
	resetReq.Header.Set("Content-Type", "application/json")

	resetRR := httptest.NewRecorder()
	router.ServeHTTP(resetRR, resetReq)

	if resetRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resetRR.Code)
	}

	sessionReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	sessionReq.AddCookie(&http.Cookie{Name: oauth2.SessionCookieName, Value: sessionID})

	sessionRR := httptest.NewRecorder()
	router.ServeHTTP(sessionRR, sessionReq)

	if sessionRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", sessionRR.Code)
	}

	var response authdto.SessionResponse
	if err := json.NewDecoder(sessionRR.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode session response: %v", err)
	}
	if response.Authenticated {
		t.Fatalf("expected reset password to revoke existing session, got %#v", response)
	}
}

func seedResetToken(
	t *testing.T,
	db *bun.DB,
	userID, rawToken string,
	expiresAt time.Time,
	usedAt *time.Time,
) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO auth_password_reset_tokens (
			password_reset_token_id, user_id, token_hash, expires_at, used_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`,
		uuid.NewString(),
		userID,
		hashResetToken(rawToken),
		expiresAt,
		usedAt,
	)
	if err != nil {
		t.Fatalf("failed to seed password reset token: %v", err)
	}
}

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
