package gamificationmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type UserAchievementModel struct {
	bun.BaseModel `bun:"table:user_achievements,alias:ua"`

	UserID         string    `bun:"user_id,pk"`
	AchievementKey string    `bun:"achievement_key,pk"`
	Level          int       `bun:"level"`
	UnlockedAt     time.Time `bun:"unlocked_at"`
}
