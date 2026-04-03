package gamification

import (
	"context"
	"errors"
	"sort"
	"strings"

	gamificationdto "github.com/master-bogdan/estimate-room-api/internal/modules/gamification/dto"
	gamificationmodels "github.com/master-bogdan/estimate-room-api/internal/modules/gamification/models"
	gamificationrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/gamification/repositories"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

const (
	AchievementSessionAdmin         = "SESSION_ADMIN"
	AchievementSessionParticipation = "SESSION_PARTICIPATION"
	AchievementTasksEstimated       = "TASKS_ESTIMATED"

	SessionRewardedEvent = "GAMIFICATION_SESSION_REWARDED"

	adminSessionXP        = 25
	participantSessionXP  = 10
	estimatedTaskXP       = 3
)

var achievementMilestones = map[string][]int{
	AchievementSessionAdmin:         {1, 5, 10},
	AchievementSessionParticipation: {1, 5, 10},
	AchievementTasksEstimated:       {1, 10, 25, 50},
}

type RewardNotifier interface {
	NotifySessionReward(ctx context.Context, reward AppliedRoomReward) error
}

type AchievementProgress struct {
	Key           string `json:"key"`
	PreviousLevel int    `json:"previousLevel"`
	CurrentLevel  int    `json:"currentLevel"`
}

type AppliedRoomReward struct {
	RoomID                     string                `json:"roomId"`
	RoomStatus                 string                `json:"roomStatus"`
	UserID                     string                `json:"userId"`
	IsAdmin                    bool                  `json:"isAdmin"`
	SessionsParticipatedDelta  int                   `json:"sessionsParticipatedDelta"`
	SessionsAdminedDelta       int                   `json:"sessionsAdminedDelta"`
	TasksEstimatedDelta        int                   `json:"tasksEstimatedDelta"`
	XPGained                   int                   `json:"xpGained"`
	PreviousXP                 int                   `json:"previousXp"`
	CurrentXP                  int                   `json:"currentXp"`
	PreviousLevel              int                   `json:"previousLevel"`
	CurrentLevel               int                   `json:"currentLevel"`
	UnlockedAchievements       []AchievementProgress `json:"unlockedAchievements"`
}

type RoomRewardService interface {
	ApplyRoomTerminalRewards(ctx context.Context, db bun.IDB, room *roomsmodels.RoomsModel) ([]AppliedRoomReward, error)
	NotifyAppliedRewards(ctx context.Context, rewards []AppliedRoomReward) error
}

type GamificationService interface {
	RoomRewardService
	GetMe(ctx context.Context, userID string) (gamificationdto.MeResponse, error)
}

type gamificationService struct {
	db          *bun.DB
	repoFactory func(db bun.IDB) gamificationrepositories.GamificationRepository
	notifier    RewardNotifier
}

func NewGamificationService(db *bun.DB, notifier RewardNotifier) GamificationService {
	return &gamificationService{
		db:          db,
		repoFactory: gamificationrepositories.NewGamificationRepository,
		notifier:    notifier,
	}
}

func (s *gamificationService) GetMe(ctx context.Context, userID string) (gamificationdto.MeResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return gamificationdto.MeResponse{}, apperrors.ErrBadRequest
	}
	if s.db == nil {
		return gamificationdto.MeResponse{}, apperrors.ErrInternal
	}

	repo := s.repoFactory(s.db)
	stats, err := repo.GetUserStats(ctx, userID)
	if err != nil {
		return gamificationdto.MeResponse{}, err
	}

	achievements, err := repo.ListUserAchievements(ctx, userID)
	if err != nil {
		return gamificationdto.MeResponse{}, err
	}

	responseAchievements := make([]gamificationdto.AchievementResponse, 0, len(achievements))
	for _, achievement := range achievements {
		if achievement == nil {
			continue
		}
		responseAchievements = append(responseAchievements, gamificationdto.AchievementResponse{
			Key:        achievement.AchievementKey,
			Level:      achievement.Level,
			UnlockedAt: achievement.UnlockedAt,
		})
	}

	sort.Slice(responseAchievements, func(i, j int) bool {
		return responseAchievements[i].Key < responseAchievements[j].Key
	})

	return gamificationdto.MeResponse{
		Stats: gamificationdto.StatsResponse{
			SessionsParticipated: stats.SessionsParticipated,
			SessionsAdmined:      stats.SessionsAdmined,
			TasksEstimated:       stats.TasksEstimated,
			XP:                   stats.XP,
			Level:                levelForXP(stats.XP),
			NextLevelXP:          nextLevelXP(stats.XP),
		},
		Achievements: responseAchievements,
	}, nil
}

func (s *gamificationService) ApplyRoomTerminalRewards(
	ctx context.Context,
	db bun.IDB,
	room *roomsmodels.RoomsModel,
) ([]AppliedRoomReward, error) {
	if room == nil || strings.TrimSpace(room.RoomID) == "" || strings.TrimSpace(room.AdminUserID) == "" {
		return nil, apperrors.ErrBadRequest
	}
	if room.Status != "FINISHED" && room.Status != "EXPIRED" {
		return nil, apperrors.ErrBadRequest
	}
	if db == nil {
		return nil, apperrors.ErrInternal
	}

	repo := s.repoFactory(db)
	candidates, err := repo.ListRoomRewardCandidates(ctx, room.RoomID, room.AdminUserID)
	if err != nil {
		return nil, err
	}

	applied := make([]AppliedRoomReward, 0, len(candidates))
	for _, candidate := range candidates {
		xpGained := rewardXP(candidate)
		rewardModel := &gamificationmodels.UserSessionRewardModel{
			RoomID:                    room.RoomID,
			UserID:                    candidate.UserID,
			IsAdmin:                   candidate.IsAdmin,
			SessionsParticipatedDelta: candidate.SessionsParticipatedDelta,
			SessionsAdminedDelta:      candidate.SessionsAdminedDelta,
			TasksEstimatedDelta:       candidate.TasksEstimatedDelta,
			XPGained:                  xpGained,
		}

		inserted, err := repo.InsertUserSessionReward(ctx, rewardModel)
		if err != nil {
			return nil, err
		}
		if !inserted {
			continue
		}

		previousStats, err := repo.GetUserStats(ctx, candidate.UserID)
		if err != nil {
			return nil, err
		}

		currentStats, err := repo.ApplyUserStatsDelta(ctx, candidate.UserID, gamificationrepositories.UserStatsDelta{
			SessionsParticipated: candidate.SessionsParticipatedDelta,
			SessionsAdmined:      candidate.SessionsAdminedDelta,
			TasksEstimated:       candidate.TasksEstimatedDelta,
			XP:                   xpGained,
		})
		if err != nil {
			return nil, err
		}

		unlockedAchievements, err := s.applyAchievements(ctx, repo, candidate.UserID, previousStats, currentStats)
		if err != nil {
			return nil, err
		}

		applied = append(applied, AppliedRoomReward{
			RoomID:                    room.RoomID,
			RoomStatus:                room.Status,
			UserID:                    candidate.UserID,
			IsAdmin:                   candidate.IsAdmin,
			SessionsParticipatedDelta: candidate.SessionsParticipatedDelta,
			SessionsAdminedDelta:      candidate.SessionsAdminedDelta,
			TasksEstimatedDelta:       candidate.TasksEstimatedDelta,
			XPGained:                  xpGained,
			PreviousXP:                previousStats.XP,
			CurrentXP:                 currentStats.XP,
			PreviousLevel:             levelForXP(previousStats.XP),
			CurrentLevel:              levelForXP(currentStats.XP),
			UnlockedAchievements:      unlockedAchievements,
		})
	}

	return applied, nil
}

func (s *gamificationService) NotifyAppliedRewards(ctx context.Context, rewards []AppliedRoomReward) error {
	if s.notifier == nil || len(rewards) == 0 {
		return nil
	}

	for _, reward := range rewards {
		if err := s.notifier.NotifySessionReward(ctx, reward); err != nil {
			return err
		}
	}

	return nil
}

func (s *gamificationService) applyAchievements(
	ctx context.Context,
	repo gamificationrepositories.GamificationRepository,
	userID string,
	previousStats, currentStats *gamificationmodels.UserStatsModel,
) ([]AchievementProgress, error) {
	updates := make([]AchievementProgress, 0, len(achievementMilestones))

	for achievementKey, milestones := range achievementMilestones {
		var previousValue int
		var currentValue int

		switch achievementKey {
		case AchievementSessionAdmin:
			previousValue = previousStats.SessionsAdmined
			currentValue = currentStats.SessionsAdmined
		case AchievementSessionParticipation:
			previousValue = previousStats.SessionsParticipated
			currentValue = currentStats.SessionsParticipated
		case AchievementTasksEstimated:
			previousValue = previousStats.TasksEstimated
			currentValue = currentStats.TasksEstimated
		default:
			continue
		}

		previousLevel := milestoneLevel(previousValue, milestones)
		currentLevel := milestoneLevel(currentValue, milestones)
		if currentLevel == 0 || currentLevel <= previousLevel {
			continue
		}

		existingAchievement, err := repo.GetUserAchievement(ctx, userID, achievementKey)
		if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
			return nil, err
		}

		storedPreviousLevel := 0
		if existingAchievement != nil {
			storedPreviousLevel = existingAchievement.Level
		}
		if storedPreviousLevel >= currentLevel {
			continue
		}

		if err := repo.SaveUserAchievement(ctx, &gamificationmodels.UserAchievementModel{
			UserID:         userID,
			AchievementKey: achievementKey,
			Level:          currentLevel,
		}); err != nil {
			return nil, err
		}

		updates = append(updates, AchievementProgress{
			Key:           achievementKey,
			PreviousLevel: storedPreviousLevel,
			CurrentLevel:  currentLevel,
		})
	}

	sort.Slice(updates, func(i, j int) bool {
		return updates[i].Key < updates[j].Key
	})

	return updates, nil
}

func rewardXP(candidate gamificationrepositories.RoomRewardCandidate) int {
	xp := candidate.TasksEstimatedDelta * estimatedTaskXP
	if candidate.IsAdmin {
		xp += adminSessionXP
		return xp
	}

	xp += participantSessionXP
	return xp
}

func milestoneLevel(value int, milestones []int) int {
	level := 0
	for idx, milestone := range milestones {
		if value >= milestone {
			level = idx + 1
		}
	}

	return level
}

func levelForXP(xp int) int {
	if xp < 0 {
		xp = 0
	}

	return (xp / 100) + 1
}

func nextLevelXP(xp int) int {
	level := levelForXP(xp)
	return level * 100
}
