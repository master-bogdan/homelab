package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/infra/email"
	authrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/auth/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	oauth2repositories "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/repositories"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/uptrace/bun"
)

type AuthModule struct {
	Controller AuthController
	Service    AuthService
}

type AuthModuleDeps struct {
	Router          chi.Router
	DB              *bun.DB
	UserService     users.UsersService
	Oauth2Service   oauth2.Oauth2Service
	SessionService  oauth2.AuthService
	FrontendBaseURL string
	EmailClient     email.Client
	Github          oauth2utils.GithubConfig
}

func NewAuthModule(deps AuthModuleDeps) *AuthModule {
	passwordResetTokenRepo := authrepositories.NewPasswordResetTokenRepository(deps.DB)
	authCodeRepo := oauth2repositories.NewOauth2AuthCodeRepository(deps.DB)
	accessTokenRepo := oauth2repositories.NewOauth2AccessTokenRepository(deps.DB)
	refreshTokenRepo := oauth2repositories.NewOauth2RefreshTokenRepository(deps.DB)
	oidcSessionRepo := oauth2repositories.NewOauth2OidcSessionRepository(deps.DB)

	service := NewAuthService(AuthServiceDeps{
		UserService:            deps.UserService,
		Oauth2Service:          deps.Oauth2Service,
		SessionService:         deps.SessionService,
		FrontendBaseURL:        deps.FrontendBaseURL,
		EmailClient:            deps.EmailClient,
		PasswordResetTokenRepo: passwordResetTokenRepo,
		AuthCodeRepo:           authCodeRepo,
		AccessTokenRepo:        accessTokenRepo,
		RefreshTokenRepo:       refreshTokenRepo,
		OidcSessionRepo:        oidcSessionRepo,
		Github:                 deps.Github,
	})
	controller := NewAuthController(service)

	deps.Router.Route("/auth", func(r chi.Router) {
		r.Post("/login", controller.Login)
		r.Post("/register", controller.Register)
		r.Post("/forgot-password", controller.ForgotPassword)
		r.Get("/reset-password/validate", controller.ValidateResetPasswordToken)
		r.Post("/reset-password", controller.ResetPassword)
		r.Post("/logout", controller.Logout)
		r.Get("/session", controller.GetSession)
		r.Get("/github/login", controller.GithubLogin)
		r.Get("/github/callback", controller.GithubCallback)
	})

	return &AuthModule{
		Controller: controller,
		Service:    service,
	}
}
