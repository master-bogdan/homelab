package notes

import (
	"log/slog"
	"net/http"

	"github.com/master-bogdan/ephermal-notes/internal/infra/db/redis"
	"github.com/redis/go-redis/v9"
)

func RouterNew(m *http.ServeMux, client *redis.Client, logger *slog.Logger) {
	repo := memory_db.NewNotesRepository(client)
	service := NewNotesService(repo)
	controller := NewNotesController(service)

	m.HandleFunc("GET /api/v1/notes/{id}", controller.GetNote)
	m.HandleFunc("POST /api/v1/notes", controller.CreateNote)
	m.HandleFunc("DELETE /api/v1/notes/{id}", controller.DeleteNote)
}
