package oauth2

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
)

type Oauth2Module struct {
	Controller Oauth2Controller
	Service    Oauth2Service
}

type Oauth2ModuleDeps struct {
	Router           chi.Router
	TokenKey         string
	Issuer           string
	ClientRepo       repositories.Oauth2ClientRepository
	AuthCodeRepo     repositories.Oauth2AuthCodeRepository
	UserRepo         repositories.Oauth2UserRepository
	OidcSessionRepo  repositories.Oauth2OidcSessionRepository
	RefreshTokenRepo repositories.Oauth2RefreshTokenRepository
	AccessTokenRepo  repositories.Oauth2AccessTokenRepository
}

func NewOauth2Module(deps Oauth2ModuleDeps) *Oauth2Module {
	svc := NewOauth2Service(
		deps.ClientRepo,
		deps.AuthCodeRepo,
		deps.UserRepo,
		deps.OidcSessionRepo,
		deps.RefreshTokenRepo,
		deps.AccessTokenRepo,
		[]byte(deps.TokenKey),
		deps.Issuer,
	)

	ctrl := NewOauth2Controller(svc)

	deps.Router.Route("/oauth2", func(r chi.Router) {
		r.Get("/authorize", ctrl.Authorize)
		r.Get("/login", ctrl.ShowLoginForm)
		r.Post("/login", ctrl.Login)
		r.Post("/token", ctrl.GetTokens)
	})

	return &Oauth2Module{
		Controller: ctrl,
		Service:    svc,
	}
}
