package roomsrepositories

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
)

type RoomsRepository interface {
	Create(model *roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error)
	FindByID(roomID string) (*roomsmodels.RoomsModel, error)
}

type roomsRepository struct {
	db *pgxpool.Pool
}

func NewRoomsRepository(db *pgxpool.Pool) *roomsRepository {
	return &roomsRepository{db: db}
}

func (r *roomsRepository) Create(model *roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error) {
	const query = `
		INSERT INTO rooms (code, name, admin_user_id, team_id, deck_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING room_id
	`

	var teamID *string
	if model.TeamID != nil && strings.TrimSpace(*model.TeamID) != "" {
		teamID = model.TeamID
	}

	var roomID string
	err := r.db.QueryRow(
		context.Background(),
		query,
		model.Code,
		model.Name,
		model.AdminUserID,
		teamID,
		string(model.DeckID),
	).Scan(&roomID)
	if err != nil {
		return nil, err
	}

	return r.FindByID(roomID)
}

func (r *roomsRepository) FindByID(roomID string) (*roomsmodels.RoomsModel, error) {
	const query = `
		SELECT room_id, code, name, admin_user_id, team_id, deck_id, status,
			allow_guests, allow_spectators, round_timer_seconds, created_at,
			last_activity_at, finished_at
		FROM rooms
		WHERE room_id = $1
	`

	var room roomsmodels.RoomsModel
	var deckID string
	err := r.db.QueryRow(context.Background(), query, roomID).Scan(
		&room.RoomID,
		&room.Code,
		&room.Name,
		&room.AdminUserID,
		&room.TeamID,
		&deckID,
		&room.Status,
		&room.AllowGuests,
		&room.AllowSpectators,
		&room.RoundTimerSeconds,
		&room.CreatedAt,
		&room.LastActivityAt,
		&room.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	room.DeckID = roomsmodels.DeckID(deckID)

	return &room, nil
}
