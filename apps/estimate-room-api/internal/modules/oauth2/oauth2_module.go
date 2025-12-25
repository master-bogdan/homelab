// package oauth2
//
// import (
// 	"log/slog"
//
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/master-bogdan/clear-cash-api/internal/infra/db/postgresql/repositories"
// )
//
// type Oauth2Module struct {
// 	Controller Oauth2Controller
// 	Service    Oauth2Service
// }
//
// func NewOauth2Module(
// 	router fiber.Router,
// 	tokenKey string,
// 	clientRepo repositories.Oauth2ClientRepository,
// 	authCodeRepo repositories.Oauth2AuthCodeRepository,
// 	userRepo repositories.Oauth2UserRepository,
// 	oidcSessionRepo repositories.Oauth2OidcSessionRepository,
// 	refreshTokenRepo repositories.Oauth2RefreshTokenRepository,
// 	accessTokenRepo repositories.Oauth2AccessTokenRepository,
// 	logger *slog.Logger,
// ) *Oauth2Module {
// 	oauthLogger := logger.With(slog.String("module", "oauth"))
//
// 	svc := NewOauth2Service(
// 		clientRepo,
// 		authCodeRepo,
// 		userRepo,
// 		oidcSessionRepo,
// 		refreshTokenRepo,
// 		accessTokenRepo,
// 		[]byte(tokenKey),
// 		oauthLogger,
// 	)
//
// 	ctrl := NewOauth2Controller(svc, oauthLogger)
//
// 	SetupOauth2RoutesV1(router, ctrl)
//
// 	return &Oauth2Module{
// 		Controller: ctrl,
// 		Service:    svc,
// 	}
// }
