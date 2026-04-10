package gamificationrepositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	gamificationmodels "github.com/master-bogdan/estimate-room-api/internal/modules/gamification/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type UserStatsDelta struct {
	SessionsParticipated int
	SessionsAdmined      int
	TasksEstimated       int
	XP                   int
}

type RoomRewardCandidate struct {
	UserID                    string `bun:"user_id"`
	IsAdmin                   bool   `bun:"is_admin"`
	SessionsParticipatedDelta int    `bun:"sessions_participated_delta"`
	SessionsAdminedDelta      int    `bun:"sessions_admined_delta"`
	TasksEstimatedDelta       int    `bun:"tasks_estimated_delta"`
}

type GamificationRepository interface {
	GetUserStats(ctx context.Context, userID string) (*gamificationmodels.UserStatsModel, error)
	ListUserAchievements(ctx context.Context, userID string) ([]*gamificationmodels.UserAchievementModel, error)
	ListRoomRewardCandidates(ctx context.Context, roomID, adminUserID string) ([]RoomRewardCandidate, error)
	InsertUserSessionReward(ctx context.Context, model *gamificationmodels.UserSessionRewardModel) (bool, error)
	ApplyUserStatsDelta(ctx context.Context, userID string, delta UserStatsDelta) (*gamificationmodels.UserStatsModel, error)
	GetUserAchievement(ctx context.Context, userID, achievementKey string) (*gamificationmodels.UserAchievementModel, error)
	SaveUserAchievement(ctx context.Context, model *gamificationmodels.UserAchievementModel) error
}

type gamificationRepository struct {
	db bun.IDB
}

func NewGamificationRepository(db bun.IDB) GamificationRepository {
	return &gamificationRepository{db: db}
}

func (r *gamificationRepository) GetUserStats(ctx context.Context, userID string) (*gamificationmodels.UserStatsModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	stats := new(gamificationmodels.UserStatsModel)
	err := r.db.NewSelect().
		Model(stats).
		Where("us.user_id = ?", userID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &gamificationmodels.UserStatsModel{UserID: userID}, nil
		}

		return nil, err
	}

	return stats, nil
}

func (r *gamificationRepository) ListUserAchievements(ctx context.Context, userID string) ([]*gamificationmodels.UserAchievementModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	achievements := make([]*gamificationmodels.UserAchievementModel, 0)
	err := r.db.NewSelect().
		Model(&achievements).
		Where("ua.user_id = ?", userID).
		OrderExpr("ua.achievement_key ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *gamificationRepository) ListRoomRewardCandidates(
	ctx context.Context,
	roomID, adminUserID string,
) ([]RoomRewardCandidate, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	candidates := make([]RoomRewardCandidate, 0)
	query := `
		WITH reward_users AS (
			SELECT ?::text AS user_id, TRUE AS is_admin
			UNION ALL
			SELECT DISTINCT rp.user_id, FALSE AS is_admin
			FROM room_participants AS rp
			WHERE rp.room_id = ?
			  AND rp.user_id IS NOT NULL
			  AND rp.user_id <> ?
		),
		estimated_task_votes AS (
			SELECT
				rp.user_id,
				COUNT(DISTINCT t.task_id)::int AS tasks_estimated_delta
			FROM votes AS v
			JOIN room_participants AS rp ON rp.room_participants_id = v.participant_id
			JOIN tasks AS t ON t.task_id = v.task_id
			WHERE rp.room_id = ?
			  AND rp.user_id IS NOT NULL
			  AND t.status = 'ESTIMATED'
			GROUP BY rp.user_id
		)
		SELECT
			ru.user_id,
			ru.is_admin,
			CASE WHEN ru.is_admin THEN 0 ELSE 1 END AS sessions_participated_delta,
			CASE WHEN ru.is_admin THEN 1 ELSE 0 END AS sessions_admined_delta,
			COALESCE(etv.tasks_estimated_delta, 0)::int AS tasks_estimated_delta
		FROM reward_users AS ru
		LEFT JOIN estimated_task_votes AS etv ON etv.user_id = ru.user_id
		ORDER BY ru.is_admin DESC, ru.user_id ASC
	`

	if err := r.db.NewRaw(query, adminUserID, roomID, adminUserID, roomID).Scan(ctx, &candidates); err != nil {
		return nil, err
	}

	return candidates, nil
}

func (r *gamificationRepository) InsertUserSessionReward(
	ctx context.Context,
	model *gamificationmodels.UserSessionRewardModel,
) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := r.db.NewInsert().
		Model(model).
		Column(
			"room_id",
			"user_id",
			"is_admin",
			"sessions_participated_delta",
			"sessions_admined_delta",
			"tasks_estimated_delta",
			"xp_gained",
		).
		On("CONFLICT (room_id, user_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *gamificationRepository) ApplyUserStatsDelta(
	ctx context.Context,
	userID string,
	delta UserStatsDelta,
) (*gamificationmodels.UserStatsModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	stats := &gamificationmodels.UserStatsModel{
		UserID:               userID,
		SessionsParticipated: delta.SessionsParticipated,
		SessionsAdmined:      delta.SessionsAdmined,
		TasksEstimated:       delta.TasksEstimated,
		XP:                   delta.XP,
	}

	_, err := r.db.NewInsert().
		Model(stats).
		Column("user_id", "sessions_participated", "sessions_admined", "tasks_estimated", "xp").
		On("CONFLICT (user_id) DO UPDATE").
		Set("sessions_participated = us.sessions_participated + EXCLUDED.sessions_participated").
		Set("sessions_admined = us.sessions_admined + EXCLUDED.sessions_admined").
		Set("tasks_estimated = us.tasks_estimated + EXCLUDED.tasks_estimated").
		Set("xp = us.xp + EXCLUDED.xp").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *gamificationRepository) GetUserAchievement(
	ctx context.Context,
	userID, achievementKey string,
) (*gamificationmodels.UserAchievementModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	achievement := new(gamificationmodels.UserAchievementModel)
	err := r.db.NewSelect().
		Model(achievement).
		Where("ua.user_id = ?", userID).
		Where("ua.achievement_key = ?", achievementKey).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return achievement, nil
}

func (r *gamificationRepository) SaveUserAchievement(
	ctx context.Context,
	model *gamificationmodels.UserAchievementModel,
) error {
	if ctx == nil {
		ctx = context.Background()
	}

	model.UnlockedAt = time.Now().UTC()
	_, err := r.db.NewInsert().
		Model(model).
		Column("user_id", "achievement_key", "level", "unlocked_at").
		On("CONFLICT (user_id, achievement_key) DO UPDATE").
		Set("level = EXCLUDED.level").
		Set("unlocked_at = EXCLUDED.unlocked_at").
		Exec(ctx)
	return err
}
