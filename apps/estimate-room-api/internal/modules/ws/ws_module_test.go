package ws

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
)

type authServiceStub struct {
	checkAuth func(*http.Request) (string, error)
	called    bool
}

func (s *authServiceStub) CheckAuth(r *http.Request) (string, error) {
	s.called = true
	if s.checkAuth != nil {
		return s.checkAuth(r)
	}
	return "", errors.New("unexpected auth call")
}

func (*authServiceStub) CreateOidcSession(*oauth2models.OidcSessionModel) (string, error) {
	return "", nil
}

func (*authServiceStub) GetOidcSessionByID(string) (*oauth2models.OidcSessionModel, error) {
	return nil, nil
}

func (*authServiceStub) CreateAccessToken(context.Context, *oauth2models.Oauth2AccessTokenModel) error {
	return nil
}

func TestResolveIdentity_RejectsAccessTokenInQuery(t *testing.T) {
	authService := &authServiceStub{
		checkAuth: func(*http.Request) (string, error) {
			t.Fatal("expected query token rejection before auth check")
			return "", nil
		},
	}

	req := httptest.NewRequest("GET", "/ws?token=query-token", nil)

	_, err := resolveIdentity(req, authService, nil)
	if !errors.Is(err, errQueryAccessTokenNotAllowed) {
		t.Fatalf("expected query token rejection, got %v", err)
	}
	if authService.called {
		t.Fatal("expected auth service not to be called")
	}
}

func TestResolveIdentity_UsesAuthorizationHeader(t *testing.T) {
	authService := &authServiceStub{
		checkAuth: func(r *http.Request) (string, error) {
			if got := r.Header.Get("Authorization"); got != "Bearer header-token" {
				t.Fatalf("expected authorization header to be preserved, got %q", got)
			}
			return "user-123", nil
		},
	}

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Authorization", "Bearer header-token")

	identity, err := resolveIdentity(req, authService, nil)
	if err != nil {
		t.Fatalf("expected header auth to succeed, got %v", err)
	}
	if identity.Type != IdentityTypeUser || identity.UserID != "user-123" {
		t.Fatalf("unexpected identity: %+v", identity)
	}
	if !authService.called {
		t.Fatal("expected auth service to be called")
	}
}
