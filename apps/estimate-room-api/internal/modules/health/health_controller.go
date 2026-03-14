package health

import (
	"net/http"
	"sync/atomic"

	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
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

// CheckReadiness godoc
// @Summary Readiness check
// @Description Checks database and Redis connectivity.
// @Tags health
// @Produce json
// @Success 200 {object} ReadinessStatus
// @Failure 503 {object} apperrors.HttpError
// @Router /api/v1/health/readyz [get]
func (c *healthController) CheckReadiness(w http.ResponseWriter, r *http.Request) {
	if c.isGracefulShutdown != nil && c.isGracefulShutdown.Load() {
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrServiceUnavailable,
			apperrors.HttpError{Detail: "graceful shutdown"},
		))
		return
	}

	status, _ := c.svc.CheckReadiness(r.Context())

	code := http.StatusOK
	if status.Status != "OK" {
		code = http.StatusServiceUnavailable
	}

	httputils.WriteResponse(w, status, httputils.WriteResponseOptions{Status: code})
}

// CheckHealth godoc
// @Summary Liveness check
// @Description Reports service liveness status.
// @Tags health
// @Produce json
// @Success 200 {object} LivenessStatus
// @Failure 503 {object} apperrors.HttpError
// @Router /api/v1/health/healthz [get]
func (c *healthController) CheckHealth(w http.ResponseWriter, r *http.Request) {
	if c.isGracefulShutdown != nil && c.isGracefulShutdown.Load() {
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrServiceUnavailable,
			apperrors.HttpError{Detail: "graceful shutdown"},
		))
		return
	}

	status, _ := c.svc.CheckHealth(r.Context())

	code := http.StatusOK
	if status.Status != "OK" {
		code = http.StatusServiceUnavailable
	}

	httputils.WriteResponse(w, status, httputils.WriteResponseOptions{Status: code})
}
