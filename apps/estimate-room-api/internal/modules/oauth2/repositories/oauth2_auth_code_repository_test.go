package oauth2repositories_test

import (
	"testing"
	"time"

	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
	oauth2repositories "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/repositories"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
)

func TestOauth2AuthCodeRepository_CreateAndFindByCode(t *testing.T) {
	db := testutils.SetupTestDB(t)
	testutils.ResetOauthTables(t, db)
	defer db.Close()

	redirectURI := "http://localhost:5173/auth/callback"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"openid", "user"})
	userID := testutils.SeedUser(t, db, "repo-auth-code@example.com", "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-auth-code-repo")

	repo := oauth2repositories.NewOauth2AuthCodeRepository(db)
	model := &oauth2models.Oauth2AuthCodeModel{
		ClientID:            clientID,
		UserID:              userID,
		OidcSessionID:       sessionID,
		Code:                "repo-auth-code",
		RedirectURI:         redirectURI,
		Scopes:              []string{"openid", "user"},
		CodeChallenge:       "challenge-auth-code",
		CodeChallengeMethod: "S256",
		ExpiresAt:           time.Now().Add(5 * time.Minute),
	}

	if err := repo.Create(model); err != nil {
		t.Fatalf("create auth code: %v", err)
	}

	saved, err := repo.FindByCode(model.Code)
	if err != nil {
		t.Fatalf("find auth code: %v", err)
	}

	if saved.AuthCodeID == "" {
		t.Fatalf("expected auth code id to be populated")
	}
	if saved.ClientID != clientID {
		t.Fatalf("expected client id %q, got %q", clientID, saved.ClientID)
	}
	if len(saved.Scopes) != 2 || saved.Scopes[0] != "openid" || saved.Scopes[1] != "user" {
		t.Fatalf("expected scopes to round-trip, got %#v", saved.Scopes)
	}
}
