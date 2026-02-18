package auth

import (
	authrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/auth/repositories"
	"github.com/uptrace/bun"
)

type AuthModule struct {
	Service AuthService
}

type AuthModuleDeps struct {
	TokenKey string
	DB       *bun.DB
}

func NewAuthModule(deps AuthModuleDeps) *AuthModule {
	accessTokenRepo := authrepositories.NewOauth2AccessTokenRepository(deps.DB)
	oidcSessionRepo := authrepositories.NewOauth2OidcSessionRepository(deps.DB)

	svc := NewAuthService(
		deps.TokenKey,
		accessTokenRepo,
		oidcSessionRepo,
	)

	return &AuthModule{
		Service: svc,
	}
}
