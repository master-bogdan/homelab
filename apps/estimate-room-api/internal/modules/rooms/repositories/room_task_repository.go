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
	FindCurrentVotingTask(roomID string) (*roomsmodels.RoomTaskModel, error)
	SetCurrentVotingTask(roomID, taskID string) (updatedTask *roomsmodels.RoomTaskModel, previousTask *roomsmodels.RoomTaskModel, err error)
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
		Column("task_id", "room_id", "title", "description", "external_key", "status", "is_active", "final_estimate_value").
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

func (r *roomTaskRepository) FindCurrentVotingTask(roomID string) (*roomsmodels.RoomTaskModel, error) {
	task := new(roomsmodels.RoomTaskModel)
	err := r.db.NewSelect().
		Model(task).
		Where("t.room_id = ?", roomID).
		Where("t.is_active = TRUE").
		OrderExpr("t.updated_at DESC").
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

func (r *roomTaskRepository) SetCurrentVotingTask(roomID, taskID string) (*roomsmodels.RoomTaskModel, *roomsmodels.RoomTaskModel, error) {
	var previousTask *roomsmodels.RoomTaskModel
	committed := false

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	targetTask := new(roomsmodels.RoomTaskModel)
	err = tx.NewSelect().
		Model(targetTask).
		Where("t.room_id = ?", roomID).
		Where("t.task_id = ?", taskID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperrors.ErrNotFound
		}
		return nil, nil, err
	}
	if targetTask.Status == "ESTIMATED" || targetTask.Status == "SKIPPED" {
		return nil, nil, apperrors.ErrBadRequest
	}

	currentVotingTask := new(roomsmodels.RoomTaskModel)
	err = tx.NewSelect().
		Model(currentVotingTask).
		Where("t.room_id = ?", roomID).
		Where("t.is_active = TRUE").
		OrderExpr("t.updated_at DESC").
		Limit(1).
		Scan(context.Background())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil, err
	}

	if currentVotingTask.TaskID != "" && currentVotingTask.TaskID != taskID {
		result, updateErr := tx.NewUpdate().
			Model((*roomsmodels.RoomTaskModel)(nil)).
			Set("status = ?", "PENDING").
			Set("is_active = FALSE").
			Set("updated_at = NOW()").
			Where("room_id = ?", roomID).
			Where("task_id = ?", currentVotingTask.TaskID).
			Exec(context.Background())
		if updateErr != nil {
			return nil, nil, updateErr
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return nil, nil, rowsErr
		}
		if rows > 0 {
			currentVotingTask.Status = "PENDING"
			currentVotingTask.IsActive = false
			previousTask = currentVotingTask
		}
	}

	result, err := tx.NewUpdate().
		Model((*roomsmodels.RoomTaskModel)(nil)).
		Set("status = ?", "VOTING").
		Set("is_active = TRUE").
		Set("updated_at = NOW()").
		Where("room_id = ?", roomID).
		Where("task_id = ?", taskID).
		Exec(context.Background())
	if err != nil {
		return nil, nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, nil, err
	}
	if rowsAffected == 0 {
		return nil, nil, apperrors.ErrNotFound
	}

	updatedTask := new(roomsmodels.RoomTaskModel)
	err = tx.NewSelect().
		Model(updatedTask).
		Where("t.room_id = ?", roomID).
		Where("t.task_id = ?", taskID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	committed = true

	return updatedTask, previousTask, nil
}

func (r *roomTaskRepository) Update(roomID string, model *roomsmodels.RoomTaskModel) (*roomsmodels.RoomTaskModel, error) {
	result, err := r.db.NewUpdate().
		Model(model).
		Column("title", "description", "external_key", "status", "is_active", "final_estimate_value").
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
