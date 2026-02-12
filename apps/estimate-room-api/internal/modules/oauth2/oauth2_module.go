// Package oauth2 is oauth2 implementation
package oauth2

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	oauth2repositories "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/repositories"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
)

type Oauth2Module struct {
	Controller Oauth2Controller
	Service    Oauth2Service
}

type Oauth2ModuleDeps struct {
	Router      chi.Router
	DB          *pgxpool.Pool
	TokenKey    string
	Issuer      string
	UserService UserService
	AuthService auth.AuthService
	Github      oauth2utils.GithubConfig
}

func NewOauth2Module(deps Oauth2ModuleDeps) *Oauth2Module {
	clientRepo := oauth2repositories.NewOauth2ClientRepository(deps.DB)
	authCodeRepo := oauth2repositories.NewOauth2AuthCodeRepository(deps.DB)
	refreshTokenRepo := oauth2repositories.NewOauth2RefreshTokenRepository(deps.DB)

	svc := NewOauth2Service(
		clientRepo,
		authCodeRepo,
		refreshTokenRepo,
		deps.UserService,
		deps.AuthService,
		[]byte(deps.TokenKey),
		deps.Issuer,
	)

	ctrl := NewOauth2Controller(svc, deps.Github)

	deps.Router.Route("/oauth2", func(r chi.Router) {
		r.Get("/authorize", ctrl.Authorize)
		r.Get("/login", ctrl.ShowLoginForm)
		r.Post("/login", ctrl.Login)
		r.Post("/token", ctrl.GetTokens)
		r.Get("/github/login", ctrl.GithubLogin)
		r.Get("/github/callback", ctrl.GithubCallback)
	})

	return &Oauth2Module{
		Controller: ctrl,
		Service:    svc,
	}
}
