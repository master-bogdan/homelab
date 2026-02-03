package auth

type AuthModule struct {
	Service AuthService
}

type AuthModuleDeps struct {
	TokenKey        string
	AccessTokenRepo AccessTokenRepository
	OidcSessionRepo OidcSessionRepository
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
