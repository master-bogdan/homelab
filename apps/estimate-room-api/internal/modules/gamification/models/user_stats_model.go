package gamificationmodels

import "github.com/uptrace/bun"

type UserStatsModel struct {
	bun.BaseModel `bun:"table:user_stats,alias:us"`

	UserID               string `bun:"user_id,pk"`
	SessionsParticipated int    `bun:"sessions_participated"`
	SessionsAdmined      int    `bun:"sessions_admined"`
	TasksEstimated       int    `bun:"tasks_estimated"`
	XP                   int    `bun:"xp"`
}
