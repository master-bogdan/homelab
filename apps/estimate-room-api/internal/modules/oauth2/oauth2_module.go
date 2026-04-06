// Package oauth2 is oauth2 implementation
package oauth2

import (
	"github.com/go-chi/chi/v5"
	oauth2repositories "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/repositories"
	"github.com/uptrace/bun"
)

type Oauth2Module struct {
	Controller         Oauth2Controller
	Service            Oauth2Service
	SessionAuthService Oauth2SessionAuthService
}

type Oauth2ModuleDeps struct {
	Router             chi.Router
	DB                 *bun.DB
	TokenKey           string
	Issuer             string
	UserService        UserService
	SessionAuthService Oauth2SessionAuthService
	FrontendBaseURL    string
	TrustProxyHeaders  bool
}

func NewOauth2Module(deps Oauth2ModuleDeps) *Oauth2Module {
	clientRepo := oauth2repositories.NewOauth2ClientRepository(deps.DB)
	authCodeRepo := oauth2repositories.NewOauth2AuthCodeRepository(deps.DB)
	refreshTokenRepo := oauth2repositories.NewOauth2RefreshTokenRepository(deps.DB)
	accessTokenRepo := oauth2repositories.NewOauth2AccessTokenRepository(deps.DB)
	oidcSessionRepo := oauth2repositories.NewOauth2OidcSessionRepository(deps.DB)

	authService := deps.SessionAuthService
	if authService == nil {
		authService = NewOauth2SessionAuthService(deps.TokenKey, accessTokenRepo, oidcSessionRepo)
	}

	svc := NewOauth2Service(
		clientRepo,
		authCodeRepo,
		refreshTokenRepo,
		deps.UserService,
		authService,
		[]byte(deps.TokenKey),
		deps.Issuer,
	)

	ctrl := NewOauth2Controller(svc, authService, deps.FrontendBaseURL, deps.TrustProxyHeaders)

	deps.Router.Route("/oauth2", func(r chi.Router) {
		r.Get("/authorize", ctrl.Authorize)
		r.Post("/token", ctrl.GetTokens)
	})

	return &Oauth2Module{
		Controller:         ctrl,
		Service:            svc,
		SessionAuthService: authService,
	}
}
