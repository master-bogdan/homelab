package tests

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
)

func TestRegisterAndAuthenticateUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)
	email := "test1@example.com"
	pass := "password123"

	_, err := svc.RegisterUser(&oauth2dto.UserDTO{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	userID, err := svc.AuthenticateUser(&oauth2dto.UserDTO{
		Email: email,
	})
	if err != nil || userID == "" {
		t.Fatalf("failed to authenticate user: %v", err)
	}
}

func TestCreateOidcSessionAndAuthCode(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)

	redirectURI := "http://localhost/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"openid", "user"})
	userID := testutils.SeedUser(t, db, "testuser2@example.com", "password123")
	nonce := "random-nonce"

	sessionID, err := svc.CreateOidcSession(&oauth2dto.CreateOidcSessionDTO{
		UserID:   userID,
		ClientID: clientID,
		Nonce:    nonce,
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	verifier := "verifier"
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	code, err := svc.CreateAuthCode(&oauth2dto.CreateOauthCodeDTO{
		ClientID:            clientID,
		UserID:              userID,
		OidcSessionID:       sessionID,
		RedirectURI:         redirectURI,
		Scopes:              "openid user",
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
	})
	if err != nil {
		t.Fatalf("failed to create auth code: %v", err)
	}

	if code == "" {
		t.Fatal("auth code is empty")
	}
}

func TestAuthorizationCodeExchange(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)

	redirectURI := "http://localhost/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"openid", "user"})
	userID := testutils.SeedUser(t, db, "testuser3@example.com", "password123")

	sessionID, err := svc.CreateOidcSession(&oauth2dto.CreateOidcSessionDTO{
		UserID:   userID,
		ClientID: clientID,
		Nonce:    "nonce-456",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	verifier := "verifier456"
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	code, err := svc.CreateAuthCode(&oauth2dto.CreateOauthCodeDTO{
		ClientID:            clientID,
		UserID:              userID,
		OidcSessionID:       sessionID,
		RedirectURI:         redirectURI,
		Scopes:              "user openid",
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
	})
	if err != nil {
		t.Fatalf("create code: %v", err)
	}

	tokens, err := svc.GetAuthorizationTokens(&oauth2dto.GetTokenDTO{
		GrantType:    "authorization_code",
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		Code:         code,
		CodeVerifier: verifier,
	})
	if err != nil {
		t.Fatalf("token exchange failed: %v", err)
	}

	if tokens.AccessToken == "" || tokens.RefreshToken == "" || tokens.IDToken == "" {
		t.Fatal("missing tokens")
	}
}

func TestAuthCodeExpires(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)

	redirectURI := "http://localhost/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"openid", "user"})
	userID := testutils.SeedUser(t, db, "testuser4@example.com", "password123")

	sessionID, err := svc.CreateOidcSession(&oauth2dto.CreateOidcSessionDTO{
		UserID:   userID,
		ClientID: clientID,
		Nonce:    "nonce-exp",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	verifier := "verifier-exp"
	hash := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	code, err := svc.CreateAuthCode(&oauth2dto.CreateOauthCodeDTO{
		ClientID:            clientID,
		UserID:              userID,
		OidcSessionID:       sessionID,
		RedirectURI:         redirectURI,
		Scopes:              "user openid",
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
	})
	if err != nil {
		t.Fatalf("create code: %v", err)
	}

	_, err = db.Exec(context.Background(), `UPDATE oauth2_auth_codes SET expires_at = NOW() - interval '1 minute' WHERE code = $1`, code)
	if err != nil {
		t.Fatalf("expire code: %v", err)
	}

	_, err = svc.GetAuthorizationTokens(&oauth2dto.GetTokenDTO{
		GrantType:    "authorization_code",
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		Code:         code,
		CodeVerifier: verifier,
	})
	if err == nil {
		t.Fatal("expected error for expired auth code")
	}
}

func TestValidateClientScope(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)

	redirectURI := "http://localhost/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"openid"})

	err := svc.ValidateClient(&oauth2dto.AuthorizeQueryDTO{
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		ResponseType:        "code",
		Scopes:              "user",
		State:               uuid.NewString(),
		CodeChallenge:       "abc",
		CodeChallengeMethod: "S256",
		Nonce:               "nonce",
	})
	if err == nil {
		t.Fatal("expected scope validation error")
	}
}
