package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
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

var _ oauth2.Oauth2SessionAuthService = (*authServiceStub)(nil)

func TestWsModule_RejectsAccessTokenInQuery(t *testing.T) {
	authService := &authServiceStub{
		checkAuth: func(*http.Request) (string, error) {
			t.Fatal("expected query token rejection before auth check")
			return "", nil
		},
	}

	router := chi.NewRouter()
	ws.NewWsModule(ws.WsModuleDeps{
		Router:      router,
		AuthService: authService,
		TokenKey:    "0123456789abcdef0123456789abcdef",
	})

	req := httptest.NewRequest(http.MethodGet, "/ws?token=query-token", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rr.Code)
	}
	if authService.called {
		t.Fatal("expected auth service not to be called")
	}
}
