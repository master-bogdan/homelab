package memory_db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type NotesModel struct {
	ID      string `redis:"id"`
	Message string `redis:"message"`
}

type NotesRepository interface {
	GetById(ID string) (*NotesModel, error)
	Create(note *NotesModel) (*NotesModel, error)
	Delete(ID string) error
}

type notesRepository struct {
	client *redis.Client
}

func NewNotesRepository(client *redis.Client) NotesRepository {
	return &notesRepository{
		client: client,
	}
}

var ctx = context.Background()

func (r *notesRepository) GetById(ID string) (*NotesModel, error) {
	key := buildKey(ID)

	val, err := r.client.GetDel(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("Note already consumed")
	}

	if err != nil {
		return nil, err
	}

	note := &NotesModel{}
	err = json.Unmarshal([]byte(val), note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (r *notesRepository) Create(note *NotesModel) (*NotesModel, error) {
	ID := uuid.New().String()
	key := buildKey(ID)

	data, err := json.Marshal(note)
	if err != nil {
		return nil, err
	}

	ttl := 15 * time.Minute
	ok, err := r.client.SetNX(ctx, key, data, ttl).Result()

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("note ID already exists")
	}

	note.ID = ID

	return note, nil
}

func (r *notesRepository) Delete(ID string) error {
	key := buildKey(ID)
	deleted, err := r.client.Del(ctx, key).Result()

	if err != nil {
		return err
	}

	if deleted == 0 {
		// Key did not exist → expired or never created
		return errors.New("note does not exist or has already expired")
	}

	return nil
}

func buildKey(ID string) string {
	prefix := "note"

	return fmt.Sprintf("%s:%s", prefix, ID)
}
