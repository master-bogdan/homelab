package health

import (
	"context"
	"time"

	healthdto "github.com/master-bogdan/estimate-room-api/internal/modules/health/dto"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type HealthService interface {
	CheckHealth(ctx context.Context) (healthdto.LivenessStatus, error)
	CheckReadiness(ctx context.Context) (healthdto.ReadinessStatus, error)
}

type healthService struct {
	db      *bun.DB
	redis   *redis.Client
	started time.Time
}

type HealthServiceDeps struct {
	DB    *bun.DB
	Redis *redis.Client
}

func NewHealthService(deps HealthServiceDeps) HealthService {
	return &healthService{
		db:      deps.DB,
		redis:   deps.Redis,
		started: time.Now(),
	}
}

func (s *healthService) CheckHealth(ctx context.Context) (healthdto.LivenessStatus, error) {
	return healthdto.LivenessStatus{
		Status: "OK",
		Uptime: time.Since(s.started),
	}, nil
}

func (s *healthService) CheckReadiness(ctx context.Context) (healthdto.ReadinessStatus, error) {
	status := healthdto.ReadinessStatus{
		Status: "OK",
		Uptime: time.Since(s.started),
		DB:     "up",
		Redis:  "up",
	}

	if s.db == nil {
		status.DB = "down"
		status.Status = "not OK"
	} else if err := s.db.PingContext(ctx); err != nil {
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
