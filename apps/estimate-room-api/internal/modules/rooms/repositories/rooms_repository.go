package roomsrepositories

import (
	"context"
	"strings"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/uptrace/bun"
)

type RoomsRepository interface {
	Create(model *roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error)
	FindByID(roomID string) (*roomsmodels.RoomsModel, error)
}

type roomsRepository struct {
	db *bun.DB
}

func NewRoomsRepository(db *bun.DB) *roomsRepository {
	return &roomsRepository{db: db}
}

func (r *roomsRepository) Create(model *roomsmodels.RoomsModel) (*roomsmodels.RoomsModel, error) {
	if model.TeamID != nil && strings.TrimSpace(*model.TeamID) == "" {
		model.TeamID = nil
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("code", "name", "admin_user_id", "team_id", "deck_id").
		Returning("*").
		Exec(context.Background())
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
				OrderExpr("t.order_index ASC").
				Relation("Votes", func(vq *bun.SelectQuery) *bun.SelectQuery {
					return vq.OrderExpr("v.created_at ASC")
				})
		}).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return room, nil
}
