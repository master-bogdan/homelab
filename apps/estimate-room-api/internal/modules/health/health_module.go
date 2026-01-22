// Package health is a health endpoints
package health

import (
	"sync/atomic"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthModule struct {
	Controller HealthController
	Service    HealthService
}

type HealthModuleDeps struct {
	Router             chi.Router
	DB                 *pgxpool.Pool
	Redis              *redis.Client
	IsGracefulShutdown *atomic.Bool
}

func NewHealthModule(deps HealthModuleDeps) *HealthModule {
	svc := NewHealthService(HealthServiceDeps{
		DB:    deps.DB,
		Redis: deps.Redis,
	})
	ctrl := NewHealthController(svc, deps.IsGracefulShutdown)

	deps.Router.Route("/health", func(r chi.Router) {
		deps.Router.Get("/healthz", ctrl.CheckHealth)
		deps.Router.Get("/readyz", ctrl.CheckReadiness)
	})

	return &HealthModule{
		Controller: ctrl,
		Service:    svc,
	}
}
