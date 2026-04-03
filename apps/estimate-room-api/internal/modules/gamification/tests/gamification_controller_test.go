package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/gamification"
	gamificationdto "github.com/master-bogdan/estimate-room-api/internal/modules/gamification/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

func setupGamificationTest(t *testing.T) (*chi.Mux, *bun.DB, gamification.GamificationService, rooms.RoomsExpiryService) {
	t.Helper()

	db := testutils.SetupTestDB(t)

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
			user_achievements,
			user_stats,
			user_session_rewards,
			invitations,
			votes,
			task_rounds,
			tasks,
			room_participants,
			rooms,
			team_members,
			teams,
			oauth2_access_tokens,
			oauth2_refresh_tokens,
			oauth2_auth_codes,
			oauth2_oidc_sessions,
			users,
			oauth2_clients
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate gamification test tables: %v", err)
	}

	router := chi.NewRouter()
	authService := oauth2.NewAuthServiceFromDB(testutils.TestTokenKey, db)
	wsService := ws.NewService(nil, "test-gamification-events")

	var gamificationModule *gamification.GamificationModule
	var roomsModule *rooms.RoomsModule

	router.Route("/api/v1", func(r chi.Router) {
		invitesModule := invites.NewInvitesModule(invites.InvitesModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
		})

		gamificationModule = gamification.NewGamificationModule(gamification.GamificationModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			WsService:   wsService,
		})

		roomsModule = rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:         r,
			DB:             db,
			WsService:      wsService,
			AuthService:    authService,
			InvitesService: invitesModule.Service,
			RewardService:  gamificationModule.Service,
		})
	})

	return router, db, gamificationModule.Service, roomsModule.ExpiryService
}

func createGamificationAccessToken(t *testing.T, db *bun.DB, email string) (string, string) {
	t.Helper()

	redirectURI := "http://localhost:4081"
	clientID := testutils.SeedClient(t, db, redirectURI, []string{"user"})
	userID := testutils.SeedUser(t, db, email, "password123")
	sessionID := testutils.SeedSession(t, db, userID, clientID, "nonce-gamification")

	svc := testutils.NewOauth2Service(db)
	tokens, err := svc.GenerateTokenPair(context.Background(), userID, clientID, sessionID, []string{"user"})
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	return tokens.AccessToken, userID
}

func seedGamificationRoom(t *testing.T, db *bun.DB, roomID, adminUserID, status string, lastActivityAt time.Time) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO rooms (room_id, code, name, admin_user_id, deck, status, last_activity_at, finished_at)
		VALUES ($1, $2, $3, $4, '{"name":"Fibonacci","kind":"FIBONACCI","values":["1","2","3","5","8"]}'::jsonb, $5, $6, CASE WHEN $5 IN ('FINISHED', 'EXPIRED') THEN $6 ELSE NULL END)
	`, roomID, "room-"+roomID[:8], "Gamification Room", adminUserID, status, lastActivityAt.UTC())
	if err != nil {
		t.Fatalf("failed to insert gamification room: %v", err)
	}
}

func seedGamificationParticipant(t *testing.T, db *bun.DB, roomID, userID, role string) string {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, user_id, role)
		VALUES ($1, $2, $3, $4)
	`, participantID, roomID, userID, role)
	if err != nil {
		t.Fatalf("failed to insert participant: %v", err)
	}

	return participantID
}

func seedGamificationGuestParticipant(t *testing.T, db *bun.DB, roomID, guestName string) string {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, guest_name, role)
		VALUES ($1, $2, $3, 'GUEST')
	`, participantID, roomID, guestName)
	if err != nil {
		t.Fatalf("failed to insert guest participant: %v", err)
	}

	return participantID
}

func seedGamificationTask(t *testing.T, db *bun.DB, roomID, status string, isActive bool, finalEstimateValue *string) string {
	t.Helper()

	taskID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO tasks (task_id, room_id, title, status, is_active, final_estimate_value)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, taskID, roomID, "Task "+taskID[:8], status, isActive, finalEstimateValue)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	return taskID
}

func seedGamificationTaskRound(t *testing.T, db *bun.DB, taskID string, roundNumber int, status string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO task_rounds (task_id, round_number, eligible_participant_ids, status)
		VALUES ($1, $2, '[]'::jsonb, $3)
	`, taskID, roundNumber, status)
	if err != nil {
		t.Fatalf("failed to insert task round: %v", err)
	}
}

func seedGamificationVote(t *testing.T, db *bun.DB, taskID, participantID string, roundNumber int, value string) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), `
		INSERT INTO votes (votes_id, task_id, participant_id, value, round_number)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.NewString(), taskID, participantID, value, roundNumber)
	if err != nil {
		t.Fatalf("failed to insert vote: %v", err)
	}
}

func loadGamificationMe(t *testing.T, router *chi.Mux, token string) gamificationdto.MeResponse {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gamification/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from gamification me, got %d: %s", rr.Code, rr.Body.String())
	}

	var response gamificationdto.MeResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode gamification me response: %v", err)
	}

	return response
}

func TestFinishRoom_AppliesRewardsAndExposesGamificationProfile(t *testing.T) {
	router, db, _, _ := setupGamificationTest(t)
	defer db.Close()

	adminToken, adminUserID := createGamificationAccessToken(t, db, "admin@example.com")
	memberToken, memberUserID := createGamificationAccessToken(t, db, "member@example.com")

	roomID := uuid.NewString()
	seedGamificationRoom(t, db, roomID, adminUserID, "ACTIVE", time.Now().UTC())

	adminParticipantID := seedGamificationParticipant(t, db, roomID, adminUserID, "ADMIN")
	memberParticipantID := seedGamificationParticipant(t, db, roomID, memberUserID, "MEMBER")
	seedGamificationGuestParticipant(t, db, roomID, "Guest")

	finalEstimate := "5"
	estimatedTaskID := seedGamificationTask(t, db, roomID, "ESTIMATED", false, &finalEstimate)
	seedGamificationTaskRound(t, db, estimatedTaskID, 1, "REVEALED")
	seedGamificationVote(t, db, estimatedTaskID, adminParticipantID, 1, "5")
	seedGamificationVote(t, db, estimatedTaskID, memberParticipantID, 1, "5")

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, strings.NewReader(`{"status":"FINISHED"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when finishing room, got %d: %s", rr.Code, rr.Body.String())
	}

	adminProfile := loadGamificationMe(t, router, adminToken)
	if adminProfile.Stats.SessionsAdmined != 1 || adminProfile.Stats.SessionsParticipated != 0 {
		t.Fatalf("expected admin stats admined=1 participated=0, got %+v", adminProfile.Stats)
	}
	if adminProfile.Stats.TasksEstimated != 1 || adminProfile.Stats.XP != 28 {
		t.Fatalf("expected admin tasksEstimated=1 xp=28, got %+v", adminProfile.Stats)
	}
	if adminProfile.Stats.Level != 1 || adminProfile.Stats.NextLevelXP != 100 {
		t.Fatalf("expected admin level 1 nextLevelXP 100, got %+v", adminProfile.Stats)
	}

	memberProfile := loadGamificationMe(t, router, memberToken)
	if memberProfile.Stats.SessionsParticipated != 1 || memberProfile.Stats.SessionsAdmined != 0 {
		t.Fatalf("expected member stats participated=1 admined=0, got %+v", memberProfile.Stats)
	}
	if memberProfile.Stats.TasksEstimated != 1 || memberProfile.Stats.XP != 13 {
		t.Fatalf("expected member tasksEstimated=1 xp=13, got %+v", memberProfile.Stats)
	}

	var userStatsCount int
	if err := db.NewSelect().
		Table("user_stats").
		ColumnExpr("COUNT(*)").
		Scan(context.Background(), &userStatsCount); err != nil {
		t.Fatalf("failed to count user_stats: %v", err)
	}
	if userStatsCount != 2 {
		t.Fatalf("expected stats rows only for registered users, got %d", userStatsCount)
	}
}

func TestApplyRoomTerminalRewards_IsIdempotentAndUpgradesAchievements(t *testing.T) {
	_, db, gamificationService, _ := setupGamificationTest(t)
	defer db.Close()

	_, adminUserID := createGamificationAccessToken(t, db, "admin@example.com")
	_, memberUserID := createGamificationAccessToken(t, db, "member@example.com")

	if _, err := db.ExecContext(context.Background(), `
		INSERT INTO user_stats (user_id, sessions_participated, sessions_admined, tasks_estimated, xp)
		VALUES ($1, 4, 0, 9, 90)
	`, memberUserID); err != nil {
		t.Fatalf("failed to seed member stats: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `
		INSERT INTO user_achievements (user_id, achievement_key, level)
		VALUES ($1, 'SESSION_PARTICIPATION', 1)
	`, memberUserID); err != nil {
		t.Fatalf("failed to seed achievement: %v", err)
	}

	room := &roomsmodels.RoomsModel{
		RoomID:      uuid.NewString(),
		AdminUserID: adminUserID,
		Status:      "FINISHED",
	}
	seedGamificationRoom(t, db, room.RoomID, adminUserID, room.Status, time.Now().UTC())

	adminParticipantID := seedGamificationParticipant(t, db, room.RoomID, adminUserID, "ADMIN")
	memberParticipantID := seedGamificationParticipant(t, db, room.RoomID, memberUserID, "MEMBER")

	finalEstimate := "8"
	taskID := seedGamificationTask(t, db, room.RoomID, "ESTIMATED", false, &finalEstimate)
	seedGamificationTaskRound(t, db, taskID, 1, "REVEALED")
	seedGamificationVote(t, db, taskID, adminParticipantID, 1, "8")
	seedGamificationVote(t, db, taskID, memberParticipantID, 1, "8")

	firstRewards, err := gamificationService.ApplyRoomTerminalRewards(context.Background(), db, room)
	if err != nil {
		t.Fatalf("failed to apply first rewards: %v", err)
	}
	if len(firstRewards) != 2 {
		t.Fatalf("expected 2 rewards on first application, got %d", len(firstRewards))
	}

	secondRewards, err := gamificationService.ApplyRoomTerminalRewards(context.Background(), db, room)
	if err != nil {
		t.Fatalf("failed to apply second rewards: %v", err)
	}
	if len(secondRewards) != 0 {
		t.Fatalf("expected no rewards on second application, got %d", len(secondRewards))
	}

	var rewardRows int
	if err := db.NewSelect().
		Table("user_session_rewards").
		ColumnExpr("COUNT(*)").
		Scan(context.Background(), &rewardRows); err != nil {
		t.Fatalf("failed to count session rewards: %v", err)
	}
	if rewardRows != 2 {
		t.Fatalf("expected 2 persisted reward rows, got %d", rewardRows)
	}

	memberProfile, err := gamificationService.GetMe(context.Background(), memberUserID)
	if err != nil {
		t.Fatalf("failed to load member profile: %v", err)
	}
	if memberProfile.Stats.SessionsParticipated != 5 {
		t.Fatalf("expected member sessionsParticipated=5, got %d", memberProfile.Stats.SessionsParticipated)
	}
	if memberProfile.Stats.TasksEstimated != 10 {
		t.Fatalf("expected member tasksEstimated=10, got %d", memberProfile.Stats.TasksEstimated)
	}
	if memberProfile.Stats.XP != 103 {
		t.Fatalf("expected member xp=103, got %d", memberProfile.Stats.XP)
	}
	if memberProfile.Stats.Level != 2 || memberProfile.Stats.NextLevelXP != 200 {
		t.Fatalf("expected member level 2 nextLevelXP 200, got %+v", memberProfile.Stats)
	}

	foundParticipationLevel2 := false
	foundTasksEstimatedLevel2 := false
	for _, achievement := range memberProfile.Achievements {
		switch achievement.Key {
		case gamification.AchievementSessionParticipation:
			foundParticipationLevel2 = achievement.Level == 2
		case gamification.AchievementTasksEstimated:
			foundTasksEstimatedLevel2 = achievement.Level == 2
		}
	}
	if !foundParticipationLevel2 {
		t.Fatal("expected SESSION_PARTICIPATION achievement level 2")
	}
	if !foundTasksEstimatedLevel2 {
		t.Fatal("expected TASKS_ESTIMATED achievement level 2")
	}
}

func TestExpiryRoom_AppliesSameRewardRules(t *testing.T) {
	_, db, gamificationService, expiryService := setupGamificationTest(t)
	defer db.Close()

	_, adminUserID := createGamificationAccessToken(t, db, "expiry-admin@example.com")
	_, memberUserID := createGamificationAccessToken(t, db, "expiry-member@example.com")

	roomID := uuid.NewString()
	staleAt := time.Now().Add(-31 * time.Minute).UTC()
	seedGamificationRoom(t, db, roomID, adminUserID, "ACTIVE", staleAt)

	seedGamificationParticipant(t, db, roomID, adminUserID, "ADMIN")
	memberParticipantID := seedGamificationParticipant(t, db, roomID, memberUserID, "MEMBER")

	finalEstimate := "3"
	taskID := seedGamificationTask(t, db, roomID, "ESTIMATED", false, &finalEstimate)
	seedGamificationTaskRound(t, db, taskID, 1, "REVEALED")
	seedGamificationVote(t, db, taskID, memberParticipantID, 1, "3")

	expiredRooms, err := expiryService.ExpireInactiveRooms(time.Now().Add(-30 * time.Minute))
	if err != nil {
		t.Fatalf("failed to expire inactive rooms: %v", err)
	}
	if len(expiredRooms) != 1 || expiredRooms[0].RoomID != roomID {
		t.Fatalf("expected room %s to expire, got %+v", roomID, expiredRooms)
	}

	memberProfile, err := gamificationService.GetMe(context.Background(), memberUserID)
	if err != nil {
		t.Fatalf("failed to load member profile: %v", err)
	}
	if memberProfile.Stats.SessionsParticipated != 1 || memberProfile.Stats.TasksEstimated != 1 || memberProfile.Stats.XP != 13 {
		t.Fatalf("expected expiry rewards to match participation rules, got %+v", memberProfile.Stats)
	}
}
