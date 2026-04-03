package gamificationdto

import "time"

type MeResponse struct {
	Stats        StatsResponse         `json:"stats"`
	Achievements []AchievementResponse `json:"achievements"`
}

type StatsResponse struct {
	SessionsParticipated int `json:"sessionsParticipated"`
	SessionsAdmined      int `json:"sessionsAdmined"`
	TasksEstimated       int `json:"tasksEstimated"`
	XP                   int `json:"xp"`
	Level                int `json:"level"`
	NextLevelXP          int `json:"nextLevelXp"`
}

type AchievementResponse struct {
	Key        string    `json:"key"`
	Level      int       `json:"level"`
	UnlockedAt time.Time `json:"unlockedAt"`
}
