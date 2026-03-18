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

func (c *historyController) ListMySessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	query, err := historydto.ParsePaginationQuery(r.URL.Query())
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

func (c *historyController) ListTeamSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	query, err := historydto.ParsePaginationQuery(r.URL.Query())
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

	c.logger.Error("request failed", logArgs...)

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

	c.logger.Error("request failed", logArgs...)

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
