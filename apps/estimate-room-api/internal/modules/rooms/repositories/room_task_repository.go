package roomsrepositories

import (
	"context"
	"database/sql"
	"errors"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type RoomTaskRepository interface {
	Create(model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error)
	FindByRoomID(roomID string) ([]*roomsmodels.RoomTaskModel, error)
	FindByID(roomID, taskID string) (*roomsmodels.RoomTaskModel, error)
	Update(roomID string, model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error)
	Delete(roomID, taskID string) error
}

type roomTaskRepository struct {
	db *bun.DB
}

func NewRoomTaskRepository(db *bun.DB) RoomTaskRepository {
	return &roomTaskRepository{db: db}
}

func (r *roomTaskRepository) Create(model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error) {
	_, err := r.db.NewInsert().
		Model(model).
		Column("task_id", "room_id", "title", "description", "external_key", "status", "final_estimate_value").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *roomTaskRepository) FindByRoomID(roomID string) ([]*roomsmodels.RoomTaskModel, error) {
	tasks := make([]*roomsmodels.RoomTaskModel, 0)
	err := r.db.NewSelect().
		Model(&tasks).
		Where("t.room_id = ?", roomID).
		OrderExpr("t.created_at ASC").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *roomTaskRepository) FindByID(roomID, taskID string) (*roomsmodels.RoomTaskModel, error) {
	task := new(roomsmodels.RoomTaskModel)
	err := r.db.NewSelect().
		Model(task).
		Where("t.room_id = ?", roomID).
		Where("t.task_id = ?", taskID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return task, nil
}

func (r *roomTaskRepository) Update(roomID string, model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error) {
	result, err := r.db.NewUpdate().
		Model(model).
		Column("title", "description", "external_key", "status", "final_estimate_value").
		Set("updated_at = NOW()").
		WherePK().
		Where("room_id = ?", roomID).
		Exec(context.Background())
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

	return r.FindByID(roomID, model.TaskID)
}

func (r *roomTaskRepository) Delete(roomID, taskID string) error {
	result, err := r.db.NewDelete().
		Model((*roomsmodels.RoomTaskModel)(nil)).
		Where("task_id = ?", taskID).
		Where("room_id = ?", roomID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
