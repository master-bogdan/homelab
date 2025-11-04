package health

import (
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func RouterNew(m *http.ServeMux, client *redis.Client, logger *slog.Logger) {
	controller := newHealthController(client, logger)

	m.HandleFunc("GET /api/v1/health/readyz", controller.getReadyz)
	m.HandleFunc("GET /api/v1/health/healthz", controller.getHealthz)
}
