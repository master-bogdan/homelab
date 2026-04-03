package gamification

import (
	stdErrors "errors"
	"log/slog"
	"net/http"

	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type GamificationController interface {
	GetMe(w http.ResponseWriter, r *http.Request)
}

type gamificationController struct {
	service     GamificationService
	authService oauth2.AuthService
	logger      *slog.Logger
}

func NewGamificationController(service GamificationService, authService oauth2.AuthService) GamificationController {
	return &gamificationController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "gamification")),
	}
}

// GetMe godoc
// @Summary Current user gamification
// @Description Returns cumulative stats and unlocked achievements for the current authenticated user.
// @Tags gamification
// @Produce json
// @Success 200 {object} gamificationdto.MeResponse
// @Failure 401 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/gamification/me [get]
func (c *gamificationController) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		switch {
		case stdErrors.Is(err, oauth2.ErrMissingToken):
			httputils.WriteResponseError(w, apperrors.CreateHttpError(
				apperrors.ErrUnauthorized,
				apperrors.HttpError{
					Detail:   "missing access token",
					Instance: r.URL.Path,
				},
			))
		default:
			httputils.WriteResponseError(w, apperrors.CreateHttpError(
				apperrors.ErrUnauthorized,
				apperrors.HttpError{
					Detail:   "invalid or expired access token",
					Instance: r.URL.Path,
				},
			))
		}
		return
	}

	response, err := c.service.GetMe(r.Context(), userID)
	if err != nil {
		c.logger.Error("failed to get gamification profile", "err", err)
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrInternal,
			apperrors.HttpError{
				Instance: r.URL.Path,
			},
		))
		return
	}

	httputils.WriteResponse(w, response)
}
