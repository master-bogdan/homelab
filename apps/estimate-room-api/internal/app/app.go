// Package app wire up application
package app

import (
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/config"
	_ "github.com/master-bogdan/estimate-room-api/docs"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/modules/health"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

type AppDeps struct {
	DB                 *pgxpool.Pool
	Redis              *redis.Client
	Cfg                *config.Config
	Router             chi.Router
	IsGracefulShutdown *atomic.Bool
	Ws                 *ws.WsServer
}

func (deps *AppDeps) SetupApp() {
	deps.Router.Use(
		logger.RequestIDMiddleware,
		middleware.RealIP,
		logger.RequestLoggerMiddleware,
		middleware.Recoverer,
		httprate.LimitByIP(100, 1*time.Minute),
	)

	deps.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	wsManager := ws.NewManager(deps.Ws, "app")

	clientRepo := repositories.NewOauth2ClientRepository(deps.DB)
	authCodeRepo := repositories.NewOauth2AuthCodeRepository(deps.DB)
	userRepo := repositories.NewUserRepository(deps.DB)
	oidcSessionRepo := repositories.NewOauth2OidcSessionRepository(deps.DB)
	refreshTokenRepo := repositories.NewOauth2RefreshTokenRepository(deps.DB)
	accessTokenRepo := repositories.NewOauth2AccessTokenRepository(deps.DB)

	deps.Router.Route("/api/v1", func(r chi.Router) {
		health.NewHealthModule(health.HealthModuleDeps{
			Router:             r,
			DB:                 deps.DB,
			Redis:              deps.Redis,
			IsGracefulShutdown: deps.IsGracefulShutdown,
		})

		rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:    r,
			WsManager: wsManager,
		})

		oauth2.NewOauth2Module(oauth2.Oauth2ModuleDeps{
			Router:           r,
			TokenKey:         deps.Cfg.Server.PasetoSymmetricKey,
			Issuer:           deps.Cfg.Server.Issuer,
			ClientRepo:       clientRepo,
			AuthCodeRepo:     authCodeRepo,
			UserRepo:         userRepo,
			OidcSessionRepo:  oidcSessionRepo,
			RefreshTokenRepo: refreshTokenRepo,
			AccessTokenRepo:  accessTokenRepo,
		})

		authModule := auth.NewAuthModule(auth.AuthModuleDeps{
			TokenKey:        deps.Cfg.Server.PasetoSymmetricKey,
			AccessTokenRepo: accessTokenRepo,
			OidcSessionRepo: oidcSessionRepo,
		})

		users.NewUsersModule(users.UsersModuleDeps{
			Router:      r,
			AuthService: authModule.Service,
			UserRepo:    userRepo,
		})
	})
}
