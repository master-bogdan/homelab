package users

import (
	stdErrors "errors"
	"log/slog"
	"net/http"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	usersdto "github.com/master-bogdan/estimate-room-api/internal/modules/users/dto"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/errors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

type UsersController interface {
	GetMe(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	service UsersService
	logger  *slog.Logger
}

func NewUsersController(service UsersService) UsersController {
	return &usersController{
		service: service,
		logger:  logger.L().With(slog.String("module", "users")),
	}
}

// GetMe godoc
// @Summary Current user
// @Description Returns the current authenticated user profile.
// @Tags users
// @Produce json
// @Success 200 {object} usersdto.UserResponse
// @Failure 401 {object} errors.Problem
// @Failure 404 {object} errors.Problem
// @Failure 500 {object} errors.Problem
// @Router /api/v1/users/me [get]
func (c *usersController) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := c.service.GetCurrentUser(r)
	if err != nil {
		switch {
		case stdErrors.Is(err, auth.ErrMissingToken):
			errors.Write(w, errors.Problem{
				Type:     "https://api.estimateroom.com/problems/unauthorized",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   "missing access token",
				Instance: r.URL.Path,
				Errors:   []errors.ErrorItem{},
			})
		case stdErrors.Is(err, ErrUnauthorized):
			errors.Write(w, errors.Problem{
				Type:     "https://api.estimateroom.com/problems/unauthorized",
				Title:    "Unauthorized",
				Status:   http.StatusUnauthorized,
				Detail:   "invalid or expired access token",
				Instance: r.URL.Path,
				Errors:   []errors.ErrorItem{},
			})
		case stdErrors.Is(err, repositories.ErrUserNotFound):
			errors.Write(w, errors.Problem{
				Type:     "https://api.estimateroom.com/problems/not-found",
				Title:    "Not Found",
				Status:   http.StatusNotFound,
				Detail:   "user not found",
				Instance: r.URL.Path,
				Errors:   []errors.ErrorItem{},
			})
		default:
			c.logger.Error("failed to get current user", "err", err)
			errors.Write(w, errors.Problem{
				Type:     "https://api.estimateroom.com/problems/internal-server-error",
				Title:    "Internal Server Error",
				Status:   http.StatusInternalServerError,
				Detail:   "internal server error",
				Instance: r.URL.Path,
				Errors:   []errors.ErrorItem{},
			})
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

	utils.WriteResponse(w, http.StatusOK, response)
}
