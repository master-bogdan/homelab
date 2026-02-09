package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type LivenessStatus struct {
	Status string        `json:"status"`
	Uptime time.Duration `json:"uptime" swaggertype:"string"`
}

type ReadinessStatus struct {
	Status string        `json:"status"`
	Uptime time.Duration `json:"uptime" swaggertype:"string"`
	DB     string        `json:"db"`
	Redis  string        `json:"redis"`
}

type HealthService interface {
	CheckHealth(ctx context.Context) (LivenessStatus, error)
	CheckReadiness(ctx context.Context) (ReadinessStatus, error)
}

type healthService struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	started time.Time
}

type HealthServiceDeps struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func NewHealthService(deps HealthServiceDeps) HealthService {
	return &healthService{
		db:      deps.DB,
		redis:   deps.Redis,
		started: time.Now(),
	}
}

func (s *healthService) CheckHealth(ctx context.Context) (LivenessStatus, error) {
	return LivenessStatus{
		Status: "OK",
		Uptime: time.Since(s.started),
	}, nil
}

func (s *healthService) CheckReadiness(ctx context.Context) (ReadinessStatus, error) {
	status := ReadinessStatus{
		Status: "OK",
		Uptime: time.Since(s.started),
		DB:     "up",
		Redis:  "up",
	}

	if s.db == nil {
		status.DB = "down"
		status.Status = "not OK"
	} else if err := s.db.Ping(ctx); err != nil {
		status.DB = "down"
		status.Status = "not OK"
	}

	if s.redis == nil {
		status.Redis = "down"
		status.Status = "not OK"
	} else if err := s.redis.Ping(ctx).Err(); err != nil {
		status.Redis = "down"
		status.Status = "not OK"
	}

	return status, nil
}
