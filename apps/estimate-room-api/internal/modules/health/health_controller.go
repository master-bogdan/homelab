package health

import (
	"net/http"
	"sync/atomic"

	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

type HealthController interface {
	CheckReadiness(w http.ResponseWriter, r *http.Request)
	CheckHealth(w http.ResponseWriter, r *http.Request)
}

type healthController struct {
	svc                HealthService
	isGracefulShutdown *atomic.Bool
}

func NewHealthController(
	healthService HealthService,
	isGracefulShutdown *atomic.Bool,
) HealthController {
	return &healthController{
		svc:                healthService,
		isGracefulShutdown: isGracefulShutdown,
	}
}

func (c *healthController) CheckReadiness(w http.ResponseWriter, r *http.Request) {
	if c.isGracefulShutdown != nil && c.isGracefulShutdown.Load() {
		utils.WriteResponseError(w, http.StatusServiceUnavailable, "graceful shutdown")
		return
	}

	status, _ := c.svc.CheckReadiness(r.Context())

	code := http.StatusOK
	if status.Status != "OK" {
		code = http.StatusServiceUnavailable
	}

	utils.WriteResponse(w, code, status)
}

func (c *healthController) CheckHealth(w http.ResponseWriter, r *http.Request) {
	if c.isGracefulShutdown != nil && c.isGracefulShutdown.Load() {
		utils.WriteResponseError(w, http.StatusServiceUnavailable, "graceful shutdown")
		return
	}

	status, _ := c.svc.CheckReadiness(r.Context())

	code := http.StatusOK
	if status.Status != "OK" {
		code = http.StatusServiceUnavailable
	}

	utils.WriteResponse(w, code, status)
}
