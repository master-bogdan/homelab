package rooms

import (
	"encoding/json"
	stdErrors "errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	roomsdto "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/dto"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
)

func (c *roomsController) CreateTask(w http.ResponseWriter, r *http.Request) {
	if !c.ensureAuthorized(w, r) {
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

	task, err := c.service.CreateTask(roomID, CreateTaskInput{
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
	if !c.ensureAuthorized(w, r) {
		return
	}

	roomID := chi.URLParam(r, "id")

	tasks, err := c.service.ListTasks(roomID)
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, tasks)
}

func (c *roomsController) GetTask(w http.ResponseWriter, r *http.Request) {
	if !c.ensureAuthorized(w, r) {
		return
	}

	roomID := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")

	task, err := c.service.GetTask(roomID, taskID)
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, task)
}

func (c *roomsController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if !c.ensureAuthorized(w, r) {
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

	task, err := c.service.UpdateTask(roomID, taskID, UpdateTaskInput{
		Title:              dto.Title,
		Description:        dto.Description,
		ExternalKey:        dto.ExternalKey,
		Status:             dto.Status,
		FinalEstimateValue: dto.FinalEstimateValue,
	})
	if err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, task)
}

func (c *roomsController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if !c.ensureAuthorized(w, r) {
		return
	}

	roomID := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")

	if err := c.service.DeleteTask(roomID, taskID); err != nil {
		c.writeTaskError(w, r, err)
		return
	}

	httputils.WriteResponse(w, map[string]bool{"ok": true})
}

func (c *roomsController) ensureAuthorized(w http.ResponseWriter, r *http.Request) bool {
	_, err := c.authService.CheckAuth(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
		return false
	}

	return true
}

func (c *roomsController) writeTaskError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case stdErrors.Is(err, apperrors.ErrNotFound):
		c.writeError(w, r, apperrors.ErrNotFound, err.Error(), err)
	case stdErrors.Is(err, apperrors.ErrBadRequest):
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}
