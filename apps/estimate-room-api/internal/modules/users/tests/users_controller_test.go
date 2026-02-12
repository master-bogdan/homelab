package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	usersdto "github.com/master-bogdan/estimate-room-api/internal/modules/users/dto"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
)

func setupUsersTest(t *testing.T) (*chi.Mux, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)

	router := chi.NewRouter()

	authModule := auth.NewAuthModule(auth.AuthModuleDeps{
		TokenKey: testutils.TestTokenKey,
		DB:       db,
	})

	router.Route("/api/v1", func(r chi.Router) {
		users.NewUsersModule(users.UsersModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authModule.Service,
		})
	})

	return router, db
}

func TestGetMe_ReturnsUser(t *testing.T) {
	router, db := setupUsersTest(t)
	defer db.Close()

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, "testuser@example.com", "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce123")

	svc := testutils.NewOauth2Service(db)
	tokens, err := svc.GenerateTokenPair(context.Background(), userID, clientID, sessionID, []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}

	var resp usersdto.UserResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, resp.ID)
	}
	if resp.Email == nil || *resp.Email != "testuser@example.com" {
		t.Fatalf("expected email testuser@example.com, got %#v", resp.Email)
	}
}

func TestGetMe_MissingToken_ReturnsUnauthorized(t *testing.T) {
	router, db := setupUsersTest(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rr.Code)
	}

	var resp apperrors.HttpError
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.Status)
	}
	if resp.Type != "https://api.estimateroom.com/problems/unauthorized" {
		t.Fatalf("expected type unauthorized, got %s", resp.Type)
	}
	if resp.Title != "Unauthorized" {
		t.Fatalf("expected title Unauthorized, got %s", resp.Title)
	}
	if resp.Errors == nil {
		t.Fatal("expected errors array")
	}
}
