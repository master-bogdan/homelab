package notes_test

import (
	"os"
	"testing"
	"time"

	memory_db "github.com/master-bogdan/ephermal-notes/internal/infra/db/redis"
	"github.com/master-bogdan/ephermal-notes/internal/modules/notes"
	"github.com/redis/go-redis/v9"
)

func setupTestRepo(t *testing.T) memory_db.NotesRepository {
	t.Helper()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to flush db: %v", err)
	}

	return memory_db.NewNotesRepository(client)
}

// -------- Positive --------

func TestCreateNote_ReturnsCreatedNote(t *testing.T) {
	repo := setupTestRepo(t)
	service := notes.NewNotesService(repo)

	note := &memory_db.NotesModel{ID: "1", Message: "hello"}
	created, err := service.CreateNote(note)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Message != "hello" {
		t.Errorf("expected hello, got %s", created.Message)
	}
}

func TestGetNote_ReturnsExistingNote(t *testing.T) {
	repo := setupTestRepo(t)
	service := notes.NewNotesService(repo)

	note := &memory_db.NotesModel{ID: "2", Message: "world"}
	_, _ = service.CreateNote(note)

	got, err := service.GetNote("2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "2" {
		t.Errorf("expected ID=2, got %s", got.ID)
	}
}

func TestDeleteNote_RemovesNote(t *testing.T) {
	repo := setupTestRepo(t)
	service := notes.NewNotesService(repo)

	note := &memory_db.NotesModel{ID: "3", Message: "temp"}
	_, _ = service.CreateNote(note)

	if err := service.DeleteNote("3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := service.GetNote("3"); err == nil {
		t.Errorf("expected error after delete, got nil")
	}
}

// -------- Negative --------

func TestGetNote_NotFound(t *testing.T) {
	repo := setupTestRepo(t)
	service := notes.NewNotesService(repo)

	if _, err := service.GetNote("does-not-exist"); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteNote_NotFound(t *testing.T) {
	repo := setupTestRepo(t)
	service := notes.NewNotesService(repo)

	if err := service.DeleteNote("does-not-exist"); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetNote_Expired(t *testing.T) {
	repo := setupTestRepo(t)
	service := notes.NewNotesService(repo)

	note := &memory_db.NotesModel{ID: "exp1", Message: "bye"}
	_, _ = service.CreateNote(note)

	time.Sleep(2 * time.Second)

	if _, err := service.GetNote("exp1"); err == nil {
		t.Errorf("expected expired error, got nil")
	}
}
