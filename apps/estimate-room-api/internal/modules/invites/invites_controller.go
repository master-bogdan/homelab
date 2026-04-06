package invites

import (
	"encoding/json"
	stdErrors "errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type InvitesController interface {
	PreviewInvitation(w http.ResponseWriter, r *http.Request)
	AcceptInvitation(w http.ResponseWriter, r *http.Request)
	DeclineInvitation(w http.ResponseWriter, r *http.Request)
	RevokeInvitation(w http.ResponseWriter, r *http.Request)
}

type invitesController struct {
	service     InvitesService
	authService oauth2.Oauth2SessionAuthService
	logger      *slog.Logger
}

func NewInvitesController(service InvitesService, authService oauth2.Oauth2SessionAuthService) InvitesController {
	return &invitesController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "invites")),
	}
}

func (c *invitesController) PreviewInvitation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	invitation, err := c.service.PreviewInvitation(r.Context(), token)
	if err != nil {
		c.writeInviteError(w, r, err)
		return
	}

	httputils.WriteResponse(w, invitesdto.NewInvitationResponse(invitation))
}

func (c *invitesController) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	userID, _ := c.optionalUserID(r)

	dto := invitesdto.AcceptInvitationDTO{}
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil && !stdErrors.Is(err, io.EOF) {
			c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
			return
		}
	}
	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	logger.FromRequest(r, c.logger).Info("accept invitation dto accepted",
		"path", r.URL.Path,
		"guest_name_provided", dto.GuestName != nil,
	)

	result, err := c.service.AcceptInvitation(r.Context(), token, userID, dto.GuestName)
	if err != nil {
		c.writeInviteError(w, r, err)
		return
	}

	if result.GuestToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     GuestAccessCookieName,
			Value:    result.GuestToken,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			Path:     "/api/v1/rooms/",
			SameSite: http.SameSiteLaxMode,
		})
	}

	if result.Room != nil {
		httputils.WriteResponse(w, map[string]any{
			"room":        result.Room,
			"participant": result.Participant,
		})
		return
	}

	httputils.WriteResponse(w, invitesdto.NewInvitationResponse(result.Invitation))
}

func (c *invitesController) DeclineInvitation(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	userID, _ := c.optionalUserID(r)

	invitation, err := c.service.DeclineInvitation(r.Context(), token, userID)
	if err != nil {
		c.writeInviteError(w, r, err)
		return
	}

	httputils.WriteResponse(w, invitesdto.NewInvitationResponse(invitation))
}

func (c *invitesController) RevokeInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	invitationID := chi.URLParam(r, "id")
	invitation, err := c.service.RevokeInvitation(r.Context(), invitationID, userID)
	if err != nil {
		c.writeInviteError(w, r, err)
		return
	}

	httputils.WriteResponse(w, invitesdto.NewInvitationResponse(invitation))
}

func (c *invitesController) writeInviteError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
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

func (c *invitesController) writeError(w http.ResponseWriter, r *http.Request, errType error, detail string, cause error) {
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

func (c *invitesController) requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return "", false
	}

	return userID, true
}

func (c *invitesController) optionalUserID(r *http.Request) (string, bool) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		return "", false
	}

	return userID, true
}
