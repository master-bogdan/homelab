// Package app wire up application
package app

import (
	"net/http"
	"strings"
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
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
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
	githubScopes := strings.Fields(deps.Cfg.Github.Scopes)

	deps.Router.Route("/api/v1", func(r chi.Router) {
		authModule := auth.NewAuthModule(auth.AuthModuleDeps{
			TokenKey:        deps.Cfg.Server.PasetoSymmetricKey,
			AccessTokenRepo: accessTokenRepo,
			OidcSessionRepo: oidcSessionRepo,
		})

		health.NewHealthModule(health.HealthModuleDeps{
			Router:             r,
			DB:                 deps.DB,
			Redis:              deps.Redis,
			IsGracefulShutdown: deps.IsGracefulShutdown,
		})

		rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:      r,
			WsManager:   wsManager,
			AuthService: authModule.Service,
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
			Github: oauth2utils.GithubConfig{
				ClientID:     deps.Cfg.Github.ClientID,
				ClientSecret: deps.Cfg.Github.ClientSecret,
				RedirectURL:  deps.Cfg.Github.RedirectURL,
				StateSecret:  deps.Cfg.Github.StateSecret,
				Scopes:       githubScopes,
			},
		})

		users.NewUsersModule(users.UsersModuleDeps{
			Router:      r,
			AuthService: authModule.Service,
			UserRepo:    userRepo,
		})

		r.Get("/ws", func(w http.ResponseWriter, req *http.Request) {
			userID, err := authModule.Service.CheckAuth(req)
			if err != nil {
				utils.WriteResponseError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			wsManager.Connect(w, req, userID)
		})
	})
}
