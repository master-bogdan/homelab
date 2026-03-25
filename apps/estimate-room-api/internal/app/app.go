// Package app wire up application
package app

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/master-bogdan/estimate-room-api/config"
	_ "github.com/master-bogdan/estimate-room-api/docs"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/modules/gamification"
	"github.com/master-bogdan/estimate-room-api/internal/modules/health"
	"github.com/master-bogdan/estimate-room-api/internal/modules/history"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/teams"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/metrics"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bun"
)

type AppDeps struct {
	DB                 *bun.DB
	Redis              *redis.Client
	Cfg                *config.Config
	Router             chi.Router
	IsGracefulShutdown *atomic.Bool
	WsServer           ws.PubSub
}

func (deps *AppDeps) SetupApp(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	if deps.Cfg == nil {
		deps.Cfg = &config.Config{}
	}

	httpRateLimitPerMinute := 100
	wsRateLimitPerMinute := 120
	if deps.Cfg.Server.HTTPRateLimitPerMinute > 0 {
		httpRateLimitPerMinute = deps.Cfg.Server.HTTPRateLimitPerMinute
	}
	if deps.Cfg.Server.WSRateLimitPerMinute > 0 {
		wsRateLimitPerMinute = deps.Cfg.Server.WSRateLimitPerMinute
	}

	deps.Router.Use(
		logger.RequestIDMiddleware,
		middleware.RealIP,
		logger.RequestLoggerMiddleware,
		middleware.Recoverer,
		httprate.LimitByIP(httpRateLimitPerMinute, 1*time.Minute),
	)

	deps.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	deps.Router.Handle("/metrics", metrics.Handler())

	githubScopes := strings.Fields(deps.Cfg.Github.Scopes)
	wsOriginPatterns := splitConfigList(deps.Cfg.Server.WebSocketAllowedOrigins)

	deps.Router.Route("/api/v1", func(r chi.Router) {
		health.NewHealthModule(health.HealthModuleDeps{
			Router:             r,
			DB:                 deps.DB,
			Redis:              deps.Redis,
			IsGracefulShutdown: deps.IsGracefulShutdown,
		})

		userRepo := usersrepositories.NewUserRepository(deps.DB)
		userService := users.NewUsersService(userRepo)

		oauth2Module := oauth2.NewOauth2Module(oauth2.Oauth2ModuleDeps{
			Router:          r,
			DB:              deps.DB,
			TokenKey:        deps.Cfg.Server.PasetoSymmetricKey,
			Issuer:          deps.Cfg.Server.Issuer,
			UserService:     userService,
			FrontendBaseURL: deps.Cfg.Frontend.BaseURL,
		})

		auth.NewAuthModule(auth.AuthModuleDeps{
			Router:         r,
			DB:             deps.DB,
			UserService:    userService,
			Oauth2Service:  oauth2Module.Service,
			SessionService: oauth2Module.AuthService,
			Github: oauth2utils.GithubConfig{
				ClientID:     deps.Cfg.Github.ClientID,
				ClientSecret: deps.Cfg.Github.ClientSecret,
				RedirectURL:  deps.Cfg.Github.RedirectURL,
				StateSecret:  deps.Cfg.Github.StateSecret,
				Scopes:       githubScopes,
			},
		})

		wsModule := ws.NewWsModule(ws.WsModuleDeps{
			Router:               r,
			AuthService:          oauth2Module.AuthService,
			TokenKey:             deps.Cfg.Server.PasetoSymmetricKey,
			Server:               deps.WsServer,
			OriginPatterns:       wsOriginPatterns,
			MessageRatePerMinute: wsRateLimitPerMinute,
		})

		users.NewUsersModule(users.UsersModuleDeps{
			Router:      r,
			DB:          deps.DB,
			AuthService: oauth2Module.AuthService,
		})

		invitesModule := invites.NewInvitesModule(invites.InvitesModuleDeps{
			Router:      r,
			DB:          deps.DB,
			AuthService: oauth2Module.AuthService,
			TokenKey:    deps.Cfg.Server.PasetoSymmetricKey,
		})

		teams.NewTeamsModule(teams.TeamsModuleDeps{
			Router:         r,
			DB:             deps.DB,
			AuthService:    oauth2Module.AuthService,
			UserService:    userService,
			InvitesService: invitesModule.Service,
		})

		gamificationModule := gamification.NewGamificationModule(gamification.GamificationModuleDeps{
			Router:      r,
			DB:          deps.DB,
			AuthService: oauth2Module.AuthService,
			WsService:   wsModule.Service,
		})

		roomsModule := rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:         r,
			DB:             deps.DB,
			WsService:      wsModule.Service,
			AuthService:    oauth2Module.AuthService,
			InvitesService: invitesModule.Service,
			RewardService:  gamificationModule.Service,
		})

		history.NewHistoryModule(history.HistoryModuleDeps{
			Router:      r,
			DB:          deps.DB,
			AuthService: oauth2Module.AuthService,
		})

		if roomsModule != nil && roomsModule.ExpiryService != nil {
			roomsModule.ExpiryService.Start(ctx)
		}
	})
}

func splitConfigList(value string) []string {
	fields := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\n' || r == '\t'
	})

	items := make([]string, 0, len(fields))
	for _, field := range fields {
		if trimmed := strings.TrimSpace(field); trimmed != "" {
			items = append(items, trimmed)
		}
	}

	return items
}
