package gamification

import (
	"context"
	"encoding/json"

	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
)

type wsRewardNotifier struct {
	wsService *ws.Service
}

type sessionRewardEventPayload struct {
	RoomID                    string                `json:"roomId"`
	RoomStatus                string                `json:"roomStatus"`
	SessionsParticipatedDelta int                   `json:"sessionsParticipatedDelta"`
	SessionsAdminedDelta      int                   `json:"sessionsAdminedDelta"`
	TasksEstimatedDelta       int                   `json:"tasksEstimatedDelta"`
	XPGained                  int                   `json:"xpGained"`
	PreviousXP                int                   `json:"previousXp"`
	CurrentXP                 int                   `json:"currentXp"`
	PreviousLevel             int                   `json:"previousLevel"`
	CurrentLevel              int                   `json:"currentLevel"`
	UnlockedAchievements      []AchievementProgress `json:"unlockedAchievements"`
}

func newWSRewardNotifier(wsService *ws.Service) RewardNotifier {
	if wsService == nil {
		return nil
	}

	return &wsRewardNotifier{wsService: wsService}
}

func (n *wsRewardNotifier) NotifySessionReward(ctx context.Context, reward AppliedRoomReward) error {
	if n == nil || n.wsService == nil {
		return nil
	}

	payload, err := json.Marshal(sessionRewardEventPayload{
		RoomID:                    reward.RoomID,
		RoomStatus:                reward.RoomStatus,
		SessionsParticipatedDelta: reward.SessionsParticipatedDelta,
		SessionsAdminedDelta:      reward.SessionsAdminedDelta,
		TasksEstimatedDelta:       reward.TasksEstimatedDelta,
		XPGained:                  reward.XPGained,
		PreviousXP:                reward.PreviousXP,
		CurrentXP:                 reward.CurrentXP,
		PreviousLevel:             reward.PreviousLevel,
		CurrentLevel:              reward.CurrentLevel,
		UnlockedAchievements:      reward.UnlockedAchievements,
	})
	if err != nil {
		return err
	}

	return n.wsService.SendToUser(reward.UserID, ws.Event{
		Type:    SessionRewardedEvent,
		UserID:  reward.UserID,
		Payload: payload,
	})
}
