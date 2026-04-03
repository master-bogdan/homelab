package tests

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupTest(t *testing.T) (*chi.Mux, *bun.DB, string, string, string, string) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	router := chi.NewRouter()
	authService := oauth2.NewAuthServiceFromDB(testutils.TestTokenKey, db)

	router.Route("/api/v1", func(r chi.Router) {
		usersModule := users.NewUsersModule(users.UsersModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
		})

		oauth2.NewOauth2Module(oauth2.Oauth2ModuleDeps{
			Router:          r,
			DB:              db,
			TokenKey:        testutils.TestTokenKey,
			Issuer:          testutils.TestIssuer,
			UserService:     usersModule.Service,
			AuthService:     authService,
			FrontendBaseURL: "http://localhost:5173",
		})
	})

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, "testuser@example.com", "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce123")

	return router, db, clientID, userID, sessionID, redirectURI
}

func TestAuthorize_WithValidSession_ReturnsCode_Integration(t *testing.T) {
	router, db, clientID, _, sessionID, redirectURI := setupTest(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
		"&redirect_uri="+url.QueryEscape(redirectURI)+
		"&response_type=code"+
		"&state=xyz123"+
		"&scopes=user"+
		"&nonce=nonce123"+
		"&code_challenge=abc123"+
		"&code_challenge_method=S256", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionID})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", rr.Code)
	}

	loc := rr.Header().Get("Location")
	if !strings.Contains(loc, "code=") {
		t.Fatalf("expected Location header to contain auth code, got: %s", loc)
	}
}

func TestAuthorize_WithoutSession_RedirectsToFrontendLogin_Integration(t *testing.T) {
	router, db, clientID, _, _, redirectURI := setupTest(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
		"&redirect_uri="+url.QueryEscape(redirectURI)+
		"&response_type=code"+
		"&state=xyz123"+
		"&scopes=user"+
		"&nonce=nonce123"+
		"&code_challenge=abc123"+
		"&code_challenge_method=S256", nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", rr.Code)
	}

	loc, err := url.Parse(rr.Header().Get("Location"))
	if err != nil {
		t.Fatalf("failed to parse redirect location: %v", err)
	}

	if loc.Scheme != "http" || loc.Host != "localhost:5173" || loc.Path != "/login" {
		t.Fatalf("expected redirect to configured frontend login, got: %s", loc.String())
	}

	continueURL := loc.Query().Get("continue")
	if !strings.Contains(continueURL, "/api/v1/oauth2/authorize") {
		t.Fatalf("expected continue url to point back to authorize, got: %s", continueURL)
	}
}

func TestAuthorize_WithSessionForDifferentNonce_RotatesSession_Integration(t *testing.T) {
	router, db, clientID, _, sessionID, redirectURI := setupTest(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
		"&redirect_uri="+url.QueryEscape(redirectURI)+
		"&response_type=code"+
		"&state=xyz123"+
		"&scopes=user"+
		"&nonce=rotated-nonce"+
		"&code_challenge=abc123"+
		"&code_challenge_method=S256", nil)
	req.AddCookie(&http.Cookie{Name: oauth2.SessionCookieName, Value: sessionID})

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", rr.Code)
	}

	if !strings.Contains(rr.Header().Get("Location"), "code=") {
		t.Fatalf("expected redirect to contain auth code, got: %s", rr.Header().Get("Location"))
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected rotated session cookie to be set")
	}
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func TestToken_AuthorizationCodeFlow_ReturnsTokens_Integration(t *testing.T) {
	router, db, clientID, _, sessionID, redirectURI := setupTest(t)
	defer db.Close()

	codeVerifier := "testverifier1234567890"
	codeChallenge := generateCodeChallenge(codeVerifier)
	state := "xyz123"

	reqLogin := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
		"&redirect_uri="+url.QueryEscape(redirectURI)+
		"&response_type=code"+
		"&state="+url.QueryEscape(state)+
		"&scopes=user"+
		"&nonce=nonce123"+
		"&code_challenge="+url.QueryEscape(codeChallenge)+
		"&code_challenge_method=S256", nil)
	reqLogin.AddCookie(&http.Cookie{Name: oauth2.SessionCookieName, Value: sessionID})

	rrLogin := httptest.NewRecorder()
	router.ServeHTTP(rrLogin, reqLogin)

	if rrLogin.Code != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", rrLogin.Code)
	}

	loc := rrLogin.Header().Get("Location")
	u, _ := url.Parse(loc)
	code := u.Query().Get("code")
	if code == "" {
		t.Fatalf("expected auth code in redirect URL")
	}

	formToken := url.Values{}
	formToken.Set("grant_type", "authorization_code")
	formToken.Set("code", code)
	formToken.Set("redirect_uri", redirectURI)
	formToken.Set("client_id", clientID)
	formToken.Set("code_verifier", codeVerifier)

	reqToken := httptest.NewRequest("POST", "/api/v1/oauth2/token", strings.NewReader(formToken.Encode()))
	reqToken.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rrToken := httptest.NewRecorder()
	router.ServeHTTP(rrToken, reqToken)

	if rrToken.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rrToken.Code)
	}

	respBody, err := io.ReadAll(rrToken.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	bodyStr := string(respBody)
	if !strings.Contains(bodyStr, "access_token") {
		t.Fatalf("expected response to contain access_token, got %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "refresh_token") {
		t.Fatalf("expected response to contain refresh_token, got %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "id_token") {
		t.Fatalf("expected response to contain id_token, got %s", bodyStr)
	}
}

func TestAuthorize_InvalidQueryParams_ReturnsBadRequest_Integration(t *testing.T) {
	router, db, _, _, _, _ := setupTest(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?invalid=1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestGetTokens_UnsupportedGrantType_ReturnsBadRequest_Integration(t *testing.T) {
	router, db, _, _, _, _ := setupTest(t)
	defer db.Close()

	form := url.Values{}
	form.Set("grant_type", "foobar")

	req := httptest.NewRequest("POST", "/api/v1/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestAuthorize_InvalidRedirectURI_ReturnsBadRequest_Integration(t *testing.T) {
	router, db, clientID, _, _, _ := setupTest(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/v1/oauth2/authorize?client_id="+clientID+
		"&redirect_uri=http://malicious.com"+
		"&response_type=code"+
		"&code_challenge=abc123"+
		"&code_challenge_method=S256"+
		"&scopes=user"+
		"&state=xyz123"+
		"&nonce=nonce123", nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
	}
}
