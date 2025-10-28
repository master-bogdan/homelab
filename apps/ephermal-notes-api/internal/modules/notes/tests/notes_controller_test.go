package notes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	memory_db "github.com/master-bogdan/ephermal-notes/internal/infra/db/redis"
	"github.com/master-bogdan/ephermal-notes/internal/modules/notes"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func setupTestController(t *testing.T) (memory_db.NotesRepository, notes.NotesService, notes.NotesController) {
	t.Helper()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   1,
	})

	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to flush DB: %v", err)
	}

	repo := memory_db.NewNotesRepository(client)
	service := notes.NewNotesService(repo)
	controller := notes.NewNotesController(service)

	return repo, service, controller
}

// ---------------- Positive Tests ----------------

func TestNoteControllerPositive(t *testing.T) {
	_, _, controller := setupTestController(t)

	t.Run("creates a note", func(t *testing.T) {
		body := []byte(`{"message":"hello world"}`)
		request := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		controller.CreateNote(response, request)

		resp := response.Result()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}

		var note memory_db.NotesModel
		if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if note.Message != "hello world" {
			t.Errorf("expected message 'hello world', got %s", note.Message)
		}
	})

	t.Run("return note by id", func(t *testing.T) {
		_ = &memory_db.NotesModel{ID: "1", Message: "test"}
		controller.CreateNote(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer([]byte(`{"message":"test"}`))))

		req := httptest.NewRequest(http.MethodGet, "/notes/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), "id", "1"))
		w := httptest.NewRecorder()

		controller.GetNote(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var got memory_db.NotesModel
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		if got.Message != "test" {
			t.Errorf("expected message 'test', got %s", got.Message)
		}
	})

	t.Run("deletes note", func(t *testing.T) {
		controller.CreateNote(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer([]byte(`{"message":"to delete"}`))))

		req := httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
		req = req.WithContext(context.WithValue(req.Context(), "id", "1"))
		w := httptest.NewRecorder()

		controller.DeleteNote(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})
}

// ---------------- Negative Tests ----------------

func TestGetNote_InvalidID(t *testing.T) {
	_, _, controller := setupTestController(t)

	t.Run("returns 404 on missing note", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/", nil)
		w := httptest.NewRecorder()

		controller.GetNote(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("return 400 because invalid body", func(t *testing.T) {
		body := []byte(`{invalid json}`)
		req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		controller.CreateNote(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}
	})
}
