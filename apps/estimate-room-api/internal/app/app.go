// Package app wire up application
package app

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/config"
	_ "github.com/master-bogdan/estimate-room-api/docs"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/modules/health"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

type AppDeps struct {
	DB                 *pgxpool.Pool
	Redis              *redis.Client
	Cfg                *config.Config
	Router             chi.Router
	IsGracefulShutdown *atomic.Bool
	WsServer           ws.PubSub
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

	githubScopes := strings.Fields(deps.Cfg.Github.Scopes)

	deps.Router.Route("/api/v1", func(r chi.Router) {
		authModule := auth.NewAuthModule(auth.AuthModuleDeps{
			TokenKey: deps.Cfg.Server.PasetoSymmetricKey,
			DB:       deps.DB,
		})

		wsModule := ws.NewWsModule(ws.WsModuleDeps{
			Router:      r,
			AuthService: authModule.Service,
			Server:      deps.WsServer,
		})

		health.NewHealthModule(health.HealthModuleDeps{
			Router:             r,
			DB:                 deps.DB,
			Redis:              deps.Redis,
			IsGracefulShutdown: deps.IsGracefulShutdown,
		})

		rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:      r,
			WsService:   wsModule.Service,
			AuthService: authModule.Service,
		})

		usersModule := users.NewUsersModule(users.UsersModuleDeps{
			Router:      r,
			DB:          deps.DB,
			AuthService: authModule.Service,
		})

		oauth2.NewOauth2Module(oauth2.Oauth2ModuleDeps{
			Router:      r,
			DB:          deps.DB,
			TokenKey:    deps.Cfg.Server.PasetoSymmetricKey,
			Issuer:      deps.Cfg.Server.Issuer,
			UserService: usersModule.Service,
			AuthService: authModule.Service,
			Github: oauth2utils.GithubConfig{
				ClientID:     deps.Cfg.Github.ClientID,
				ClientSecret: deps.Cfg.Github.ClientSecret,
				RedirectURL:  deps.Cfg.Github.RedirectURL,
				StateSecret:  deps.Cfg.Github.StateSecret,
				Scopes:       githubScopes,
			},
		})
	})
}
