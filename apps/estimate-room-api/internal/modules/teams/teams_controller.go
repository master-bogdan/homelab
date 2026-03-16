package teams

import (
	"encoding/json"
	stdErrors "errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	teamsdto "github.com/master-bogdan/estimate-room-api/internal/modules/teams/dto"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type TeamsController interface {
	CreateTeam(w http.ResponseWriter, r *http.Request)
	ListTeams(w http.ResponseWriter, r *http.Request)
	GetTeam(w http.ResponseWriter, r *http.Request)
	RemoveMember(w http.ResponseWriter, r *http.Request)
}

type teamsController struct {
	service     TeamsService
	authService oauth2.AuthService
	logger      *slog.Logger
}

func NewTeamsController(service TeamsService, authService oauth2.AuthService) TeamsController {
	return &teamsController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "teams")),
	}
}

func (c *teamsController) CreateTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	dto := teamsdto.CreateTeamDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	team, err := c.service.CreateTeam(r.Context(), dto.Name, userID)
	if err != nil {
		c.writeTeamError(w, r, err)
		return
	}

	httputils.WriteResponse(w, teamsdto.NewTeamDetailResponse(team))
}

func (c *teamsController) ListTeams(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	teams, err := c.service.ListTeams(userID)
	if err != nil {
		c.writeTeamError(w, r, err)
		return
	}

	response := make([]teamsdto.TeamSummaryResponse, 0, len(teams))
	for _, team := range teams {
		response = append(response, teamsdto.NewTeamSummaryResponse(team))
	}

	httputils.WriteResponse(w, response)
}

func (c *teamsController) GetTeam(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	teamID := chi.URLParam(r, "id")
	team, err := c.service.GetTeam(teamID, userID)
	if err != nil {
		c.writeTeamError(w, r, err)
		return
	}

	httputils.WriteResponse(w, teamsdto.NewTeamDetailResponse(team))
}

func (c *teamsController) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	teamID := chi.URLParam(r, "id")
	targetUserID := chi.URLParam(r, "userId")

	if err := c.service.RemoveMember(teamID, userID, targetUserID); err != nil {
		c.writeTeamError(w, r, err)
		return
	}

	httputils.WriteResponse(w, map[string]bool{"ok": true})
}

func (c *teamsController) writeTeamError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case stdErrors.Is(err, apperrors.ErrBadRequest):
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrUnauthorized):
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrForbidden):
		c.writeError(w, r, apperrors.ErrForbidden, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrNotFound):
		c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}

func (c *teamsController) writeError(w http.ResponseWriter, r *http.Request, errType error, detail string, cause error) {
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

func (c *teamsController) requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return "", false
	}

	return userID, true
}
