package tests

// import (
// 	"crypto/sha256"
// 	"database/sql"
// 	"encoding/base64"
// 	"log/slog"
// 	"os"
// 	"testing"
//
// 	"github.com/google/uuid"
// 	"github.com/master-bogdan/clear-cash-api/internal/infra/db/postgresql/repositories"
// 	"github.com/master-bogdan/clear-cash-api/internal/modules/oauth2"
// 	oauth2_dto "github.com/master-bogdan/clear-cash-api/internal/modules/oauth2/dto"
// 	test_utils "github.com/master-bogdan/clear-cash-api/internal/pkg/test"
// )
//
// var service oauth2.Oauth2Service
// var db *sql.DB
//
// func TestMain(m *testing.M) {
// 	app, database, cfg, err := test_utils.SetupTestApp()
// 	if err != nil {
// 		panic("failed to setup test app: " + err.Error())
// 	}
//
// 	db = database
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
// 	service = oauth2.NewOauth2Module(
// 		v1,
// 		cfg.Token.TokenKey,
// 		clientRepo,
// 		authCodeRepo,
// 		userRepo,
// 		oidcSessionRepo,
// 		refreshTokenRepo,
// 		accessTokenRepo,
// 		slog.Default(),
// 	).Service
//
// 	code := m.Run()
//
// 	_ = db.Close()
// 	os.Exit(code)
// }
//
// func TestRegisterAndAuthenticateUser(t *testing.T) {
// 	email := "test1@example.com"
// 	pass := "password123"
//
// 	_, err := service.RegisterUser(&oauth2_dto.UserDTO{
// 		Email:    email,
// 		Password: pass,
// 	})
// 	if err != nil {
// 		t.Fatalf("failed to register user: %v", err)
// 	}
//
// 	userID, err := service.AuthenticateUser(&oauth2_dto.UserDTO{
// 		Email: email,
// 	})
// 	if err != nil || userID == "" {
// 		t.Fatalf("failed to authenticate user: %v", err)
// 	}
// }
//
// func TestCreateOidcSessionAndAuthCode(t *testing.T) {
// 	userID := uuid.New().String()
// 	clientID := uuid.New().String()
// 	nonce := "random-nonce"
//
// 	// Create fake client
// 	_, err := db.Exec(`INSERT INTO oauth2_clients (client_id, client_secret, redirect_uris, grant_types, response_types, scopes, client_name, client_type, created_at)
// 	VALUES ($1, '', ARRAY['http://localhost/callback'], ARRAY['authorization_code'], ARRAY['code'], ARRAY['openid'], 'Test Client', 'public', NOW())`,
// 		clientID)
// 	if err != nil {
// 		t.Fatalf("failed to create test client: %v", err)
// 	}
//
// 	// Insert test user
// 	_, err = db.Exec(`
// 		INSERT INTO oauth2_users (user_id, email, password_hash)
// 		VALUES ($1, 'testuser2@example.com', 'fakehash')
// 		ON CONFLICT DO NOTHING;
// 	`, userID)
// 	if err != nil {
// 		t.Fatalf("failed to insert user: %v", err)
// 	}
//
// 	// Create session
// 	sessionID, err := service.CreateOidcSession(&oauth2_dto.CreateOidcSessionDTO{
// 		UserID:   userID,
// 		ClientID: clientID,
// 		Nonce:    nonce,
// 	})
// 	if err != nil {
// 		t.Fatalf("failed to create session: %v", err)
// 	}
//
// 	// Create code
// 	verifier := "verifier"
// 	hash := sha256.Sum256([]byte(verifier))
// 	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
//
// 	code, err := service.CreateAuthCode(&oauth2_dto.CreateOauthCodeDTO{
// 		ClientID:            clientID,
// 		UserID:              userID,
// 		OidcSessionID:       sessionID,
// 		RedirectURI:         "http://localhost/callback",
// 		Scopes:              "openid user",
// 		CodeChallenge:       challenge,
// 		CodeChallengeMethod: "S256",
// 	})
// 	if err != nil {
// 		t.Fatalf("failed to create auth code: %v", err)
// 	}
//
// 	if code == "" {
// 		t.Fatal("auth code is empty")
// 	}
// }
//
// func TestAuthorizationCodeExchange(t *testing.T) {
// 	userID := uuid.New().String()
// 	clientID := uuid.New().String()
// 	redirectURI := "http://localhost/callback"
// 	verifier := "verifier456"
// 	hash := sha256.Sum256([]byte(verifier))
// 	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
//
// 	// Insert test client
// 	_, err := db.Exec(`INSERT INTO oauth2_clients (client_id, client_secret, redirect_uris, grant_types, response_types, scopes, client_name, client_type, created_at)
// 	VALUES ($1, '', ARRAY[$2], ARRAY['authorization_code'], ARRAY['code'], ARRAY['openid'], 'Test Client 2', 'public', NOW())`,
// 		clientID, redirectURI)
// 	if err != nil {
// 		t.Fatalf("insert client: %v", err)
// 	}
//
// 	// Insert test user
// 	_, err = db.Exec(`
// 		INSERT INTO oauth2_users (user_id, email, password_hash)
// 		VALUES ($1, 'testuser3@example.com', 'fakehash')
// 		ON CONFLICT DO NOTHING;
// 	`, userID)
// 	if err != nil {
// 		t.Fatalf("failed to insert user: %v", err)
// 	}
//
// 	sessionID, err := service.CreateOidcSession(&oauth2_dto.CreateOidcSessionDTO{
// 		UserID:   userID,
// 		ClientID: clientID,
// 		Nonce:    "nonce-456",
// 	})
// 	if err != nil {
// 		t.Fatalf("create session: %v", err)
// 	}
//
// 	code, err := service.CreateAuthCode(&oauth2_dto.CreateOauthCodeDTO{
// 		ClientID:            clientID,
// 		UserID:              userID,
// 		OidcSessionID:       sessionID,
// 		RedirectURI:         redirectURI,
// 		Scopes:              "user openid",
// 		CodeChallenge:       challenge,
// 		CodeChallengeMethod: "S256",
// 	})
// 	if err != nil {
// 		t.Fatalf("create code: %v", err)
// 	}
//
// 	tokens, err := service.GetAuthorizationTokens(&oauth2_dto.GetTokenDTO{
// 		GrantType:    "authorization_code",
// 		ClientID:     clientID,
// 		Code:         code,
// 		CodeVerifier: verifier,
// 	})
// 	if err != nil {
// 		t.Fatalf("token exchange failed: %v", err)
// 	}
//
// 	if tokens.AccessToken == "" || tokens.RefreshToken == "" || tokens.IDToken == "" {
// 		t.Fatal("missing tokens")
// 	}
// }
