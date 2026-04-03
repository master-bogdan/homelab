package rooms

import (
	"encoding/json"
	stdErrors "errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	roomsdto "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/dto"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsController interface {
	CreateRoom(w http.ResponseWriter, r *http.Request)
	GetRoom(w http.ResponseWriter, r *http.Request)
	UpdateRoom(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	ListTasks(w http.ResponseWriter, r *http.Request)
	GetTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
}

type roomsController struct {
	service       RoomsService
	taskService   RoomsTaskService
	inviteService invites.InvitesService
	authService   oauth2.AuthService
	logger        *slog.Logger
}

func NewRoomsController(
	service RoomsService,
	taskService RoomsTaskService,
	inviteService invites.InvitesService,
	authService oauth2.AuthService,
) RoomsController {
	return &roomsController{
		service:       service,
		taskService:   taskService,
		inviteService: inviteService,
		authService:   authService,
		logger:        logger.L().With(slog.String("controller", "rooms")),
	}
}

func (c *roomsController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return
	}

	dto := roomsdto.CreateRoomDTO{}
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	err = dto.Validate()
	if err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	deck := roomsmodels.RoomDeck{}
	if dto.Deck != nil {
		deck = roomsmodels.RoomDeck{
			Name:   dto.Deck.Name,
			Kind:   dto.Deck.Kind,
			Values: dto.Deck.Values,
		}
	}

	createdRoom, err := c.service.CreateRoom(r.Context(), CreateRoomInput{
		Name:            dto.Name,
		Deck:            deck,
		AdminUserID:     userID,
		InviteTeamID:    stringPointerOrNil(dto.InviteTeamID),
		InviteEmails:    dto.InviteEmails,
		CreateShareLink: dto.CreateShareLink,
	})
	if err != nil {
		switch {
		case stdErrors.Is(err, apperrors.ErrForbidden):
			c.writeError(w, r, apperrors.ErrForbidden, err.Error(), err)
		case stdErrors.Is(err, apperrors.ErrNotFound):
			c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
		case stdErrors.Is(err, apperrors.ErrBadRequest):
			c.writeError(w, r, apperrors.ErrBadRequest, "failed to create room", err)
		default:
			c.writeError(w, r, apperrors.ErrInternal, "failed to create room", err)
		}
		return
	}

	response := roomsdto.CreateRoomResponse{
		Room:              createdRoom.Room,
		EmailInvites:      make([]invitesdto.InvitationWithTokenResponse, 0, len(createdRoom.EmailInvitations)),
		SkippedRecipients: make([]roomsdto.CreateRoomSkippedRecipientResponse, 0, len(createdRoom.SkippedRecipients)),
	}

	for _, invite := range createdRoom.EmailInvitations {
		response.EmailInvites = append(response.EmailInvites, invitesdto.NewInvitationWithTokenResponse(invite.Invitation, invite.Token))
	}

	if createdRoom.ShareLink != nil {
		shareLink := invitesdto.NewInvitationWithTokenResponse(createdRoom.ShareLink.Invitation, createdRoom.ShareLink.Token)
		response.ShareLink = &shareLink
		response.InviteToken = createdRoom.ShareLink.Token
	}

	for _, skipped := range createdRoom.SkippedRecipients {
		response.SkippedRecipients = append(response.SkippedRecipients, roomsdto.CreateRoomSkippedRecipientResponse{
			UserID: skipped.UserID,
			Email:  skipped.Email,
			Reason: skipped.Reason,
		})
	}

	httputils.WriteResponse(w, response)
}

func (c *roomsController) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	if err := c.ensureRoomReadable(r, roomID); err != nil {
		c.writeRoomError(w, r, err)
		return
	}

	room, err := c.service.GetRoom(roomID)
	if err != nil {
		c.writeRoomError(w, r, err)
		return
	}

	httputils.WriteResponse(w, room)
}

func (c *roomsController) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")

	dto := roomsdto.UpdateRoomDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	room, err := c.service.UpdateRoom(roomID, userID, UpdateRoomInput{
		Name:   dto.Name,
		Status: dto.Status,
	})
	if err != nil {
		c.writeRoomError(w, r, err)
		return
	}

	if dto.Status != nil && *dto.Status == "FINISHED" {
		// TODO: emit websocket event when room is finished.
	}

	httputils.WriteResponse(w, room)
}

func (c *roomsController) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")

	dto := roomsdto.CreateRoomTaskDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	if strings.TrimSpace(dto.Title) == "" {
		c.writeError(w, r, apperrors.ErrBadRequest, "title is required", nil)
		return
	}

	task, err := c.taskService.CreateTask(roomID, userID, CreateTaskInput{
		Title:       dto.Title,
		Description: dto.Description,
		ExternalKey: dto.ExternalKey,
	})
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, task)
}

func (c *roomsController) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")

	tasks, err := c.taskService.ListTasks(roomID, userID)
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, tasks)
}

func (c *roomsController) GetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")

	task, err := c.taskService.GetTask(roomID, taskID, userID)
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, task)
}

func (c *roomsController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")

	dto := roomsdto.UpdateRoomTaskDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	if dto.Title != nil && strings.TrimSpace(*dto.Title) == "" {
		c.writeError(w, r, apperrors.ErrBadRequest, "title is required", nil)
		return
	}

	task, err := c.taskService.UpdateTask(roomID, taskID, userID, UpdateTaskInput{
		Title:              dto.Title,
		Description:        dto.Description,
		ExternalKey:        dto.ExternalKey,
		Status:             dto.Status,
		IsActive:           dto.IsActive,
		FinalEstimateValue: dto.FinalEstimateValue,
	})
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, task)
}

func (c *roomsController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := c.requireUserID(w, r)
	if !ok {
		return
	}

	roomID := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")

	if err := c.taskService.DeleteTask(roomID, taskID, userID); err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, map[string]bool{"ok": true})
}

func (c *roomsController) writeError(w http.ResponseWriter, r *http.Request, errType error, detail string, cause error) {
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

func (c *roomsController) writeRoomError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case stdErrors.Is(err, apperrors.ErrBadRequest):
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrUnauthorized):
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrNotFound):
		c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrConflict):
		c.writeError(w, r, apperrors.ErrConflict, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrForbidden):
		c.writeError(w, r, apperrors.ErrForbidden, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}

func (c *roomsController) writeTaskError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case stdErrors.Is(err, apperrors.ErrNotFound):
		c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrBadRequest):
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrForbidden):
		c.writeError(w, r, apperrors.ErrForbidden, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}

func (c *roomsController) requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return "", false
	}

	return userID, true
}

func (c *roomsController) optionalUserID(r *http.Request) (string, bool) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		return "", false
	}

	return userID, true
}

func (c *roomsController) ensureRoomReadable(r *http.Request, roomID string) error {
	if userID, ok := c.optionalUserID(r); ok {
		return c.service.ValidateUserRoomAccess(roomID, userID)
	}

	cookie, err := r.Cookie(invites.GuestAccessCookieName)
	if err != nil {
		return apperrors.ErrUnauthorized
	}

	_, err = c.inviteService.ValidateGuestRoomAccess(roomID, cookie.Value)
	return err
}

func stringPointerOrNil(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
