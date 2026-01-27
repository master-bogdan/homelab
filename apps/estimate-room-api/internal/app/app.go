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
	"github.com/master-bogdan/estimate-room-api/internal/modules/health"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
	"github.com/redis/go-redis/v9"
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

	wsManager := ws.NewManager(deps.Ws, "app")

	health.NewHealthModule(health.HealthModuleDeps{
		Router:             deps.Router,
		DB:                 deps.DB,
		Redis:              deps.Redis,
		IsGracefulShutdown: deps.IsGracefulShutdown,
	})

	rooms.NewRoomsModule(rooms.RoomsModuleDeps{
		Router:    deps.Router,
		WsManager: wsManager,
	})
}
