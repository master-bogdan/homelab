package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	authrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/auth/repositories"
)

type AuthModule struct {
	Service AuthService
}

type AuthModuleDeps struct {
	TokenKey string
	DB       *pgxpool.Pool
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
