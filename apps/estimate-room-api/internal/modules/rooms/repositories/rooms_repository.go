package roomsrepositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type RoomsRepository interface {
	Create(ctx context.Context, model *roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error)
	FindByID(roomID string) (*roomsmodels.RoomsModel, error)
	Update(roomID string, input UpdateRoomFields) (*roomsmodels.RoomsModel, error)
	TouchActivity(roomID string) error
	ExpireInactiveRooms(cutoff time.Time) ([]*roomsmodels.RoomsModel, error)
}

type roomsRepository struct {
	db *bun.DB
}

type UpdateRoomFields struct {
	Name   *string
	Status *string
}

func NewRoomsRepository(db *bun.DB) *roomsRepository {
	return &roomsRepository{db: db}
}

func (r *roomsRepository) Create(ctx context.Context, model *roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if model.TeamID != nil && strings.TrimSpace(*model.TeamID) == "" {
		model.TeamID = nil
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("code", "name", "admin_user_id", "team_id", "deck").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *roomsRepository) FindByID(roomID string) (*roomsmodels.RoomsModel, error) {
	room := new(roomsmodels.RoomsModel)
	err := r.db.NewSelect().
		Model(room).
		Where("r.room_id = ?", roomID).
		Relation("Participants", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.
				Where("rp.left_at IS NULL").
				OrderExpr("rp.joined_at ASC").
				Relation("User").
				Relation("Votes", func(vq *bun.SelectQuery) *bun.SelectQuery {
					return vq.OrderExpr("v.created_at ASC")
				})
		}).
		Relation("Tasks", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.
				OrderExpr("t.created_at ASC").
				Relation("Votes", func(vq *bun.SelectQuery) *bun.SelectQuery {
					return vq.OrderExpr("v.created_at ASC")
				})
		}).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return room, nil
}

func (r *roomsRepository) Update(roomID string, input UpdateRoomFields) (*roomsmodels.RoomsModel, error) {
	query := r.db.NewUpdate().
		Model((*roomsmodels.RoomsModel)(nil)).
		Set("last_activity_at = NOW()").
		Where("room_id = ?", roomID)

	if input.Name != nil {
		query = query.Set("name = ?", *input.Name)
	}

	if input.Status != nil {
		query = query.Set("status = ?", *input.Status)
		if *input.Status == "FINISHED" || *input.Status == "EXPIRED" {
			query = query.Set("finished_at = COALESCE(finished_at, NOW())")
		} else {
			query = query.Set("finished_at = NULL")
		}
	}

	result, err := query.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, apperrors.ErrNotFound
	}

	return r.FindByID(roomID)
}

func (r *roomsRepository) TouchActivity(roomID string) error {
	_, err := r.db.NewUpdate().
		Model((*roomsmodels.RoomsModel)(nil)).
		Set("last_activity_at = NOW()").
		Where("room_id = ?", roomID).
		Where("status = ?", "ACTIVE").
		Exec(context.Background())

	return err
}

func (r *roomsRepository) ExpireInactiveRooms(cutoff time.Time) ([]*roomsmodels.RoomsModel, error) {
	rooms := make([]*roomsmodels.RoomsModel, 0)
	err := r.db.NewRaw(`
		UPDATE rooms
		SET status = 'EXPIRED',
		    finished_at = COALESCE(finished_at, NOW())
		WHERE status = 'ACTIVE'
		  AND last_activity_at <= ?
		RETURNING *
	`, cutoff).
		Scan(context.Background(), &rooms)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rooms, nil
		}
		return nil, err
	}

	return rooms, nil
}
