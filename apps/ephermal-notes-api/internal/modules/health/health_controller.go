package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type HealthController interface {
	getReadyz(w http.ResponseWriter, r *http.Request)
	getHealthz(w http.ResponseWriter, r *http.Request)
}

type healthController struct {
	client *redis.Client
	logger *slog.Logger
}

func newHealthController(client *redis.Client, logger *slog.Logger) HealthController {
	return &healthController{
		client: client,
		logger: logger,
	}
}

func (c *healthController) getReadyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := c.client.Ping(ctx).Err(); err != nil {
		c.logger.Error("readiness check failed", "component", "redis", "error", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"unhealthy","component":"redis"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}

func (c *healthController) getHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
