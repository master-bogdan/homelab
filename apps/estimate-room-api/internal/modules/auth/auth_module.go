package auth

import "github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"

type AuthModule struct {
	Service AuthService
}

type AuthModuleDeps struct {
	TokenKey        string
	AccessTokenRepo repositories.Oauth2AccessTokenRepository
	OidcSessionRepo repositories.Oauth2OidcSessionRepository
}

func NewAuthModule(deps AuthModuleDeps) *AuthModule {
	svc := NewAuthService(
		deps.TokenKey,
		deps.AccessTokenRepo,
		deps.OidcSessionRepo,
	)

	return &AuthModule{
		Service: svc,
	}
}
