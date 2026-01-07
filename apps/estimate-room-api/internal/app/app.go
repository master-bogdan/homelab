// Package app wire up application
package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/redis/go-redis/v9"
)

type AppDeps struct {
	DB      *pgxpool.Pool
	CacheDB *redis.Client
	Cfg     *config.Config
	Router  chi.Router
}

func (deps *AppDeps) SetupApp() {
}
