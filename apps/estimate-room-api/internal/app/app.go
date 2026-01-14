// Package app wire up application
package app

import (
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/master-bogdan/estimate-room-api/internal/modules/health"
	"github.com/redis/go-redis/v9"
)

type AppDeps struct {
	DB                 *pgxpool.Pool
	Redis              *redis.Client
	Cfg                *config.Config
	Router             chi.Router
	IsGracefulShutdown *atomic.Bool
}

func (deps *AppDeps) SetupApp() {
	health.NewHealthModule(health.HealthModuleDeps{
		Router:             deps.Router,
		DB:                 deps.DB,
		Redis:              deps.Redis,
		IsGracefulShutdown: deps.IsGracefulShutdown,
	})
}
