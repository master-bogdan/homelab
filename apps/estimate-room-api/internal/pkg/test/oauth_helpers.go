package testutils

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

const (
	TestTokenKey = "0123456789abcdef0123456789abcdef"
	TestIssuer   = "http://localhost:8000"
)

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL or DATABASE_URL is required")
		return nil
	}

	if err := postgresql.MigrateUp(dbURL); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	db, err := postgresql.Connect(dbURL)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func ResetOauthTables(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	_, err := db.Exec(context.Background(), `
		TRUNCATE TABLE
			oauth2_access_tokens,
			oauth2_refresh_tokens,
			oauth2_auth_codes,
			oauth2_oidc_sessions,
			users,
			oauth2_clients
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate oauth tables: %v", err)
	}
}

func NewOauth2Service(db *pgxpool.Pool) oauth2.Oauth2Service {
	clientRepo := repositories.NewOauth2ClientRepository(db)
	authCodeRepo := repositories.NewOauth2AuthCodeRepository(db)
	userRepo := repositories.NewUserRepository(db)
	oidcSessionRepo := repositories.NewOauth2OidcSessionRepository(db)
	refreshTokenRepo := repositories.NewOauth2RefreshTokenRepository(db)
	accessTokenRepo := repositories.NewOauth2AccessTokenRepository(db)

	return oauth2.NewOauth2Service(
		clientRepo,
		authCodeRepo,
		userRepo,
		oidcSessionRepo,
		refreshTokenRepo,
		accessTokenRepo,
		[]byte(TestTokenKey),
		TestIssuer,
	)
}

func SeedClient(t *testing.T, db *pgxpool.Pool, redirectURI string, scopes []string) string {
	t.Helper()

	clientID := uuid.NewString()
	_, err := db.Exec(context.Background(), `
		INSERT INTO oauth2_clients (
			client_id, client_secret, redirect_uris, grant_types, response_types,
			scopes, client_name, client_type, created_at
		)
		VALUES ($1, '', ARRAY[$2], ARRAY['authorization_code'], ARRAY['code'], $3, 'Test Client', 'public', NOW())
	`, clientID, redirectURI, scopes)
	if err != nil {
		t.Fatalf("failed to insert client: %v", err)
	}

	return clientID
}

func SeedUser(t *testing.T, db *pgxpool.Pool, email, password string) string {
	t.Helper()

	userID := uuid.NewString()
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_, err = db.Exec(context.Background(), `
		INSERT INTO users (user_id, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, userID, email, passwordHash)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	return userID
}

func SeedSession(t *testing.T, db *pgxpool.Pool, userID, clientID, nonce string) string {
	t.Helper()

	sessionID := uuid.NewString()
	_, err := db.Exec(context.Background(), `
		INSERT INTO oauth2_oidc_sessions (oidc_session_id, user_id, client_id, nonce)
		VALUES ($1, $2, $3, $4)
	`, sessionID, userID, clientID, nonce)
	if err != nil {
		t.Fatalf("failed to insert session: %v", err)
	}

	return sessionID
}
