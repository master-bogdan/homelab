package users

import (
	stdErrors "errors"
	"log/slog"
	"net/http"

	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	usersdto "github.com/master-bogdan/estimate-room-api/internal/modules/users/dto"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type UsersController interface {
	GetMe(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	service     UsersService
	authService auth.AuthService
	logger      *slog.Logger
}

func NewUsersController(service UsersService, authService auth.AuthService) UsersController {
	return &usersController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "users")),
	}
}

// GetMe godoc
// @Summary Current user
// @Description Returns the current authenticated user profile.
// @Tags users
// @Produce json
// @Success 200 {object} usersdto.UserResponse
// @Failure 401 {object} apperrors.HttpError
// @Failure 404 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/users/me [get]
func (c *usersController) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		switch {
		case stdErrors.Is(err, auth.ErrMissingToken):
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

	user, err := c.service.GetCurrentUser(userID)
	if err != nil {
		switch {
		case stdErrors.Is(err, apperrors.ErrUserNotFound):
			httputils.WriteResponseError(w, apperrors.CreateHttpError(
				apperrors.ErrUserNotFound,
				apperrors.HttpError{
					Instance: r.URL.Path,
				},
			))
		default:
			c.logger.Error("failed to get current user", "err", err)
			httputils.WriteResponseError(w, apperrors.CreateHttpError(
				apperrors.ErrInternal,
				apperrors.HttpError{
					Instance: r.URL.Path,
				},
			))
		}
		return
	}

	response := usersdto.UserResponse{
		ID:          user.UserID,
		Email:       user.Email,
		GithubID:    user.GithubID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
		DeletedAt:   user.DeletedAt,
	}

	httputils.WriteResponse(w, response)
}
