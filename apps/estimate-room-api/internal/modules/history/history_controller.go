package history

import (
	stdErrors "errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	historydto "github.com/master-bogdan/estimate-room-api/internal/modules/history/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type HistoryController interface {
	ListMySessions(w http.ResponseWriter, r *http.Request)
	ListTeamSessions(w http.ResponseWriter, r *http.Request)
	GetRoomSummary(w http.ResponseWriter, r *http.Request)
}

type historyController struct {
	service     HistoryService
	authService oauth2.AuthService
	logger      *slog.Logger
}

func NewHistoryController(service HistoryService, authService oauth2.AuthService) HistoryController {
	return &historyController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "history")),
	}
}

// ListMySessions godoc
// @Summary Current user sessions
// @Description Returns paginated room history for the current authenticated user.
// @Tags history
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Status filter" Enums(ALL,ACTIVE,FINISHED,EXPIRED)
// @Param role query string false "Role filter" Enums(ALL,ADMIN,PARTICIPANT)
// @Success 200 {object} historydto.SessionListResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 401 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/history/me/sessions [get]
func (c *historyController) ListMySessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	query, err := historydto.ParseMySessionsQuery(r.URL.Query())
	if err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	response, err := c.service.ListMySessions(r.Context(), userID, query)
	if err != nil {
		c.writeHistoryError(w, r, err)
		return
	}

	httputils.WriteResponse(w, response)
}

// ListTeamSessions godoc
// @Summary Team sessions
// @Description Returns paginated room history for a team. Team owner access is required.
// @Tags history
// @Produce json
// @Param id path string true "Team ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param status query string false "Status filter" Enums(ALL,ACTIVE,FINISHED,EXPIRED)
// @Success 200 {object} historydto.SessionListResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 401 {object} apperrors.HttpError
// @Failure 403 {object} apperrors.HttpError
// @Failure 404 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/history/teams/{id}/sessions [get]
func (c *historyController) ListTeamSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	query, err := historydto.ParseTeamSessionsQuery(r.URL.Query())
	if err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	teamID := chi.URLParam(r, "id")
	response, err := c.service.ListTeamSessions(r.Context(), teamID, userID, query)
	if err != nil {
		c.writeHistoryError(w, r, err)
		return
	}

	httputils.WriteResponse(w, response)
}

// GetRoomSummary godoc
// @Summary Room summary
// @Description Returns aggregated room history, participants, tasks, rounds, and revealed votes.
// @Tags history
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} historydto.RoomSummaryResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 401 {object} apperrors.HttpError
// @Failure 403 {object} apperrors.HttpError
// @Failure 404 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/history/rooms/{id}/summary [get]
func (c *historyController) GetRoomSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")
	response, err := c.service.GetRoomSummary(r.Context(), roomID, userID)
	if err != nil {
		c.writeHistoryError(w, r, err)
		return
	}

	httputils.WriteResponse(w, response)
}

func (c *historyController) writeHistoryError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case stdErrors.Is(err, errHistoryNotImplemented):
		c.writeNotImplemented(w, r, "history endpoint not implemented yet", err)
	case stdErrors.Is(err, apperrors.ErrBadRequest):
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrUnauthorized):
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrForbidden):
		c.writeError(w, r, apperrors.ErrForbidden, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrNotFound):
		c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrConflict):
		c.writeError(w, r, apperrors.ErrConflict, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}

func (c *historyController) writeNotImplemented(w http.ResponseWriter, r *http.Request, detail string, cause error) {
	logArgs := []any{
		"path", r.URL.Path,
		"status", http.StatusNotImplemented,
	}
	if detail != "" {
		logArgs = append(logArgs, "detail", detail)
	}
	if cause != nil {
		logArgs = append(logArgs, "err", cause)
	}

	logger.FromRequest(r, c.logger).Error("request failed", logArgs...)

	httputils.WriteResponseError(w, apperrors.HttpError{
		Type:     "https://api.estimateroom.com/problems/not-implemented",
		Title:    "Not Implemented",
		Status:   http.StatusNotImplemented,
		Detail:   detail,
		Instance: r.URL.Path,
		Errors:   []apperrors.ErrorItem{},
	})
}

func (c *historyController) writeError(w http.ResponseWriter, r *http.Request, errType error, detail string, cause error) {
	logArgs := []any{
		"path", r.URL.Path,
		"type", errType.Error(),
	}
	if detail != "" {
		logArgs = append(logArgs, "detail", detail)
	}
	if cause != nil {
		logArgs = append(logArgs, "err", cause)
	}

	logger.FromRequest(r, c.logger).Error("request failed", logArgs...)

	httputils.WriteResponseError(w, apperrors.CreateHttpError(
		errType,
		apperrors.HttpError{
			Detail:   detail,
			Instance: r.URL.Path,
		},
	))
}

func (c *historyController) requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return "", false
	}

	return userID, true
}
