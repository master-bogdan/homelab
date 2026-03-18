package gamificationmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type UserSessionRewardModel struct {
	bun.BaseModel `bun:"table:user_session_rewards,alias:usr"`

	RoomID                     string    `bun:"room_id,pk"`
	UserID                     string    `bun:"user_id,pk"`
	IsAdmin                    bool      `bun:"is_admin"`
	SessionsParticipatedDelta  int       `bun:"sessions_participated_delta"`
	SessionsAdminedDelta       int       `bun:"sessions_admined_delta"`
	TasksEstimatedDelta        int       `bun:"tasks_estimated_delta"`
	XPGained                   int       `bun:"xp_gained"`
	CreatedAt                  time.Time `bun:"created_at"`
}
