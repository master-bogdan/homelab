package roomsrepositories

import (
	"context"
	"database/sql"
	"errors"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

func (r *roomsRepository) RoomExists(roomID string) (bool, error) {
	return r.db.NewSelect().
		Model((*roomsmodels.RoomsModel)(nil)).
		Where("r.room_id = ?", roomID).
		Exists(context.Background())
}

func (r *roomsRepository) CreateTask(model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error) {
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

func (r *roomsRepository) FindTasksByRoomID(roomID string) ([]*roomsmodels.RoomTaskModel, error) {
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

func (r *roomsRepository) FindTaskByID(roomID, taskID string) (*roomsmodels.RoomTaskModel, error) {
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

func (r *roomsRepository) UpdateTask(roomID string, model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error) {
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

	return r.FindTaskByID(roomID, model.TaskID)
}

func (r *roomsRepository) DeleteTask(roomID, taskID string) error {
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
