package rooms

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	roomsdto "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/dto"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type RoomsController interface {
	CreateRoom(w http.ResponseWriter, r *http.Request)
}

type roomsController struct {
	service     RoomsService
	authService auth.AuthService
	logger      *slog.Logger
}

func NewRoomsController(service RoomsService, authService auth.AuthService) RoomsController {
	return &roomsController{
		service:     service,
		authService: authService,
		logger:      logger.L().With(slog.String("controller", "rooms")),
	}
}

func (c *roomsController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.CheckAuth(r)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrMissingToken):
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

	dto := roomsdto.CreateRoomDTO{}
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrBadRequest,
			apperrors.HttpError{
				Detail:   err.Error(),
				Instance: r.URL.Path,
			},
		))
		return
	}

	err = dto.Validate()
	if err != nil {
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrBadRequest,
			apperrors.HttpError{
				Detail:   err.Error(),
				Instance: r.URL.Path,
			},
		))
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
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrBadRequest,
			apperrors.HttpError{
				Detail:   err.Error(),
				Instance: r.URL.Path,
			},
		))
		return
	}

	httputils.WriteResponse(w, createdRoom)
}
