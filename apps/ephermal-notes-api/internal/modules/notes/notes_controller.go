package notes

import (
	"encoding/json"
	"net/http"
	"strings"

	memory_db "github.com/master-bogdan/ephermal-notes/internal/infra/db/redis"
)

type NotesController interface {
	GetNote(w http.ResponseWriter, r *http.Request)
	CreateNote(w http.ResponseWriter, r *http.Request)
	DeleteNote(w http.ResponseWriter, r *http.Request)
}

type notesController struct {
	service NotesService
}

func NewNotesController(service NotesService) NotesController {
	return &notesController{
		service: service,
	}
}

// GetNote godoc
// @Summary Get a note by ID
// @Description Retrieve a single note using its ID
// @Tags notes
// @Accept  json
// @Produce  json
// @Param id path string true "Note ID"
// @Success 200 {object} memory_db.NotesModel
// @Failure 400 {string} string "Invalid note ID"
// @Failure 404 {string} string "Note not found"
// @Router /notes/{id} [get]
func (c *notesController) GetNote(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("id")
	if noteID == "" {
		WriteJSONError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	note, err := c.service.GetNote(noteID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "consumed") {
			WriteJSONError(w, http.StatusNotFound, "note not found")
			return
		}
		WriteJSONError(w, http.StatusInternalServerError, "failed to get note")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

type CreateNoteDTO struct {
	Message string `json:"message"`
}

// CreateNote godoc
// @Summary Create a new note
// @Description Create a note with a message
// @Tags notes
// @Accept  json
// @Produce  json
// @Param note body notes.CreateNoteDTO true "Note input"
// @Success 201 {object} memory_db.NotesModel
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to create note"
// @Router /notes [post]
func (c *notesController) CreateNote(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	noteDTO := &CreateNoteDTO{}

	err := json.NewDecoder(r.Body).Decode(noteDTO)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(noteDTO.Message) == "" {
		WriteJSONError(w, http.StatusBadRequest, "Message is required")
		return
	}

	note := &memory_db.NotesModel{
		Message: noteDTO.Message,
	}

	createdNote, err := c.service.CreateNote(note)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdNote)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// DeleteNote godoc
// @Summary Delete a note by ID
// @Description Delete an existing note
// @Tags notes
// @Accept  json
// @Produce  json
// @Param id path string true "Note ID"
// @Success 200 {string} string "Deleted successfully"
// @Failure 400 {string} string "Invalid note ID"
// @Failure 500 {string} string "Failed to delete note"
// @Router /notes/{id} [delete]
func (c *notesController) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("id")
	if noteID == "" {
		WriteJSONError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	err := c.service.DeleteNote(noteID)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
