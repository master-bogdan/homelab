package tests

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/google/uuid"
	oauth2module "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func TestRegisterAndAuthenticateUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)
	email := "test1@example.com"
	pass := "password123"

	_, err := svc.RegisterUser(&oauth2dto.Oauth2UserDTO{
		Email:    email,
		Password: pass,
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	userID, err := svc.AuthenticateUser(&oauth2dto.Oauth2UserDTO{
		Email:    email,
		Password: pass,
	})
	if err != nil || userID == "" {
		t.Fatalf("failed to authenticate user: %v", err)
	}
}

func TestAuthenticateUser_RejectsSoftDeletedUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)
	email := "deleted-auth@example.com"
	pass := "password123"
	userID := testutils.SeedUser(t, db, email, pass)
	softDeleteUser(t, db, userID)

	_, err := svc.AuthenticateUser(&oauth2dto.Oauth2UserDTO{
		Email:    email,
		Password: pass,
	})
	if !errors.Is(err, oauth2module.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials for deleted user, got %v", err)
	}
}

func TestGetOrCreateUserFromGithub_RejectsSoftDeletedGithubID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)
	userID := testutils.SeedUser(t, db, "deleted-github-id@example.com", "password123")

	_, err := db.ExecContext(context.Background(), `
		UPDATE users
		SET github_id = $1, deleted_at = NOW()
		WHERE user_id = $2
	`, "deleted-github-id", userID)
	if err != nil {
		t.Fatalf("failed to soft-delete github user: %v", err)
	}

	_, err = svc.GetOrCreateUserFromGithub(oauth2utils.GithubProfile{
		ID:          "deleted-github-id",
		DisplayName: "Deleted Github User",
	})
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatalf("expected user not found for deleted github user, got %v", err)
	}

	if got := countUsers(t, db); got != 1 {
		t.Fatalf("expected deleted github login to avoid creating users, got %d", got)
	}
}

func TestGetOrCreateUserFromGithub_RejectsSoftDeletedEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	svc := testutils.NewOauth2Service(db)
	email := "deleted-github-email@example.com"
	userID := testutils.SeedUser(t, db, email, "password123")
	softDeleteUser(t, db, userID)

	_, err := svc.GetOrCreateUserFromGithub(oauth2utils.GithubProfile{
		ID:          "new-github-id",
		Email:       &email,
		DisplayName: "Deleted Email User",
	})
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Fatalf("expected user not found for deleted github email match, got %v", err)
	}

	if got := countUsers(t, db); got != 1 {
		t.Fatalf("expected deleted github email match to avoid creating users, got %d", got)
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

	sessionID, err := svc.CreateOidcSession(&oauth2dto.Oauth2CreateOidcSessionDTO{
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

	code, err := svc.CreateAuthCode(&oauth2dto.Oauth2CreateAuthCodeDTO{
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

	sessionID, err := svc.CreateOidcSession(&oauth2dto.Oauth2CreateOidcSessionDTO{
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

	code, err := svc.CreateAuthCode(&oauth2dto.Oauth2CreateAuthCodeDTO{
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

	tokens, err := svc.GetAuthorizationTokens(context.Background(), &oauth2dto.Oauth2GetTokenDTO{
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

	sessionID, err := svc.CreateOidcSession(&oauth2dto.Oauth2CreateOidcSessionDTO{
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

	code, err := svc.CreateAuthCode(&oauth2dto.Oauth2CreateAuthCodeDTO{
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

	_, err = db.ExecContext(context.Background(), `UPDATE oauth2_auth_codes SET expires_at = NOW() - interval '1 minute' WHERE code = $1`, code)
	if err != nil {
		t.Fatalf("expire code: %v", err)
	}

	_, err = svc.GetAuthorizationTokens(context.Background(), &oauth2dto.Oauth2GetTokenDTO{
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

	err := svc.ValidateClient(&oauth2dto.Oauth2AuthorizeQueryDTO{
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

func softDeleteUser(t *testing.T, db *bun.DB, userID string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		UPDATE users
		SET deleted_at = NOW()
		WHERE user_id = $1
	`, userID)
	if err != nil {
		t.Fatalf("failed to soft-delete user: %v", err)
	}
}

func countUsers(t *testing.T, db *bun.DB) int {
	t.Helper()

	var count int
	err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count users: %v", err)
	}

	return count
}
