package roomsrepositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/uptrace/bun"
)

type RoomVoteRepository interface {
	Upsert(taskID, participantID string, roundNumber int, value string) (*roomsmodels.RoomVoteModel, error)
	ListByTaskAndRound(taskID string, roundNumber int) ([]*roomsmodels.RoomVoteModel, error)
	CountDistinctParticipantsByTaskAndRound(taskID string, roundNumber int) (int, error)
}

type roomVoteRepository struct {
	db *bun.DB
}

func NewRoomVoteRepository(db *bun.DB) RoomVoteRepository {
	return &roomVoteRepository{db: db}
}

func (r *roomVoteRepository) Upsert(taskID, participantID string, roundNumber int, value string) (*roomsmodels.RoomVoteModel, error) {
	model := &roomsmodels.RoomVoteModel{
		VoteID:        uuid.NewString(),
		TaskID:        taskID,
		ParticipantID: participantID,
		RoundNumber:   roundNumber,
		Value:         value,
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("votes_id", "task_id", "participant_id", "value", "round_number").
		On("CONFLICT (task_id, participant_id, round_number) DO UPDATE").
		Set("value = EXCLUDED.value").
		Set("created_at = NOW()").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	updated := new(roomsmodels.RoomVoteModel)
	err = r.db.NewSelect().
		Model(updated).
		Where("v.task_id = ?", taskID).
		Where("v.participant_id = ?", participantID).
		Where("v.round_number = ?", roundNumber).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return updated, nil
}

func (r *roomVoteRepository) ListByTaskAndRound(taskID string, roundNumber int) ([]*roomsmodels.RoomVoteModel, error) {
	votes := make([]*roomsmodels.RoomVoteModel, 0)
	err := r.db.NewSelect().
		Model(&votes).
		Where("v.task_id = ?", taskID).
		Where("v.round_number = ?", roundNumber).
		OrderExpr("v.created_at ASC").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return votes, nil
}

func (r *roomVoteRepository) CountDistinctParticipantsByTaskAndRound(taskID string, roundNumber int) (int, error) {
	var count int
	err := r.db.NewSelect().
		TableExpr("votes AS v").
		ColumnExpr("COUNT(DISTINCT v.participant_id)").
		Where("v.task_id = ?", taskID).
		Where("v.round_number = ?", roundNumber).
		Scan(context.Background(), &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
