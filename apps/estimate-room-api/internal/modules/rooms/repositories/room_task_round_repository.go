package roomsrepositories

import (
	"context"
	"database/sql"
	"errors"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type RoomTaskRoundRepository interface {
	GetCurrent(taskID string) (*roomsmodels.RoomTaskRoundModel, error)
	GetOrCreateCurrent(taskID string) (*roomsmodels.RoomTaskRoundModel, error)
	Advance(taskID string) (*roomsmodels.RoomTaskRoundModel, error)
	MarkRevealed(taskID string, roundNumber int) (*roomsmodels.RoomTaskRoundModel, error)
}

type roomTaskRoundRepository struct {
	db *bun.DB
}

func NewRoomTaskRoundRepository(db *bun.DB) RoomTaskRoundRepository {
	return &roomTaskRoundRepository{db: db}
}

func (r *roomTaskRoundRepository) GetCurrent(taskID string) (*roomsmodels.RoomTaskRoundModel, error) {
	model := new(roomsmodels.RoomTaskRoundModel)
	err := r.db.NewSelect().
		Model(model).
		Where("tr.task_id = ?", taskID).
		OrderExpr("tr.round_number DESC").
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return model, nil
}

func (r *roomTaskRoundRepository) GetOrCreateCurrent(taskID string) (*roomsmodels.RoomTaskRoundModel, error) {
	current, err := r.GetCurrent(taskID)
	if err == nil {
		return current, nil
	}
	if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	model := &roomsmodels.RoomTaskRoundModel{
		TaskID:      taskID,
		RoundNumber: 1,
		IsRevealed:  false,
	}

	_, err = r.db.NewInsert().
		Model(model).
		Column("task_id", "round_number", "is_revealed").
		On("CONFLICT (task_id, round_number) DO NOTHING").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return r.GetCurrent(taskID)
}

func (r *roomTaskRoundRepository) Advance(taskID string) (*roomsmodels.RoomTaskRoundModel, error) {
	current, err := r.GetOrCreateCurrent(taskID)
	if err != nil {
		return nil, err
	}

	next := &roomsmodels.RoomTaskRoundModel{
		TaskID:      taskID,
		RoundNumber: current.RoundNumber + 1,
		IsRevealed:  false,
	}

	_, err = r.db.NewInsert().
		Model(next).
		Column("task_id", "round_number", "is_revealed").
		On("CONFLICT (task_id, round_number) DO NOTHING").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return r.GetCurrent(taskID)
}

func (r *roomTaskRoundRepository) MarkRevealed(taskID string, roundNumber int) (*roomsmodels.RoomTaskRoundModel, error) {
	result, err := r.db.NewUpdate().
		Model((*roomsmodels.RoomTaskRoundModel)(nil)).
		Set("is_revealed = TRUE").
		Set("revealed_at = NOW()").
		Set("updated_at = NOW()").
		Where("task_id = ?", taskID).
		Where("round_number = ?", roundNumber).
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

	model := new(roomsmodels.RoomTaskRoundModel)
	err = r.db.NewSelect().
		Model(model).
		Where("tr.task_id = ?", taskID).
		Where("tr.round_number = ?", roundNumber).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return model, nil
}
