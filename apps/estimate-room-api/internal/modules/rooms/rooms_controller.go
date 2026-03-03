package rooms

import (
	"encoding/json"
	stdErrors "errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
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
	service     RoomsService
	authService oauth2.AuthService
	logger      *slog.Logger
}

func NewRoomsController(service RoomsService, authService oauth2.AuthService) RoomsController {
	return &roomsController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "rooms")),
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

	if strings.TrimSpace(dto.DeckID) == "" {
		dto.DeckID = string(roomsmodels.DeckIDFibonacci)
	}

	var teamID *string
	if trimmedTeamID := strings.TrimSpace(dto.TeamID); trimmedTeamID != "" {
		teamID = &trimmedTeamID
	}

	room := roomsmodels.RoomsModel{
		Name:        dto.Name,
		TeamID:      teamID,
		DeckID:      roomsmodels.DeckID(dto.DeckID),
		AdminUserID: userID,
	}

	createdRoom, err := c.service.CreateRoom(room)
	if err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	httputils.WriteResponse(w, createdRoom)
}

func (c *roomsController) GetRoom(w http.ResponseWriter, r *http.Request) {
	_, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return
	}

	roomID := chi.URLParam(r, "id")

	room, err := c.service.GetRoom(roomID)
	if err != nil {
		c.writeRoomError(w, r, err)
		return
	}

	httputils.WriteResponse(w, room)
}

func (c *roomsController) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
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
		Name:              dto.Name,
		Status:            dto.Status,
		AllowGuests:       dto.AllowGuests,
		AllowSpectators:   dto.AllowSpectators,
		RoundTimerSeconds: dto.RoundTimerSeconds,
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
	case stdErrors.Is(err, apperrors.ErrNotFound):
		c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrForbidden):
		c.writeError(w, r, apperrors.ErrForbidden, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}
