package roomsrepositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type RoomParticipantRepository interface {
	FindActiveByUserID(roomID, userID string) (*roomsmodels.RoomParticipantModel, error)
	FindActiveByGuestName(roomID, guestName string) (*roomsmodels.RoomParticipantModel, error)
	FindActiveByID(roomID, participantID string) (*roomsmodels.RoomParticipantModel, error)
	ListActiveByRoom(roomID string) ([]*roomsmodels.RoomParticipantModel, error)
	CountActiveByRoom(roomID string) (int, error)
	Create(model *roomsmodels.RoomParticipantModel) (*roomsmodels.RoomParticipantModel, error)
}

type roomParticipantRepository struct {
	db bun.IDB
}

func NewRoomParticipantRepository(db bun.IDB) RoomParticipantRepository {
	return &roomParticipantRepository{db: db}
}

func (r *roomParticipantRepository) FindActiveByUserID(roomID, userID string) (*roomsmodels.RoomParticipantModel, error) {
	participant := new(roomsmodels.RoomParticipantModel)
	err := r.db.NewSelect().
		Model(participant).
		Where("rp.room_id = ?", roomID).
		Where("rp.user_id = ?", userID).
		Where("rp.left_at IS NULL").
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return participant, nil
}

func (r *roomParticipantRepository) FindActiveByGuestName(roomID, guestName string) (*roomsmodels.RoomParticipantModel, error) {
	trimmedGuestName := strings.TrimSpace(guestName)
	participant := new(roomsmodels.RoomParticipantModel)
	err := r.db.NewSelect().
		Model(participant).
		Where("rp.room_id = ?", roomID).
		Where("rp.guest_name = ?", trimmedGuestName).
		Where("rp.left_at IS NULL").
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return participant, nil
}

func (r *roomParticipantRepository) FindActiveByID(roomID, participantID string) (*roomsmodels.RoomParticipantModel, error) {
	participant := new(roomsmodels.RoomParticipantModel)
	err := r.db.NewSelect().
		Model(participant).
		Where("rp.room_id = ?", roomID).
		Where("rp.room_participants_id = ?", participantID).
		Where("rp.left_at IS NULL").
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return participant, nil
}

func (r *roomParticipantRepository) ListActiveByRoom(roomID string) ([]*roomsmodels.RoomParticipantModel, error) {
	participants := make([]*roomsmodels.RoomParticipantModel, 0)
	err := r.db.NewSelect().
		Model(&participants).
		Where("rp.room_id = ?", roomID).
		Where("rp.left_at IS NULL").
		OrderExpr("rp.joined_at ASC").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (r *roomParticipantRepository) CountActiveByRoom(roomID string) (int, error) {
	var count int
	err := r.db.NewSelect().
		TableExpr("room_participants AS rp").
		ColumnExpr("COUNT(*)").
		Where("rp.room_id = ?", roomID).
		Where("rp.left_at IS NULL").
		Scan(context.Background(), &count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *roomParticipantRepository) Create(model *roomsmodels.RoomParticipantModel) (*roomsmodels.RoomParticipantModel, error) {
	_, err := r.db.NewInsert().
		Model(model).
		Column("room_participants_id", "room_id", "user_id", "guest_name", "role").
		Returning("*").
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return model, nil
}
