package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/gamification"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

const gamificationTestWSOrigin = "http://frontend.test"

type gamificationTestPubSub struct {
	mu            sync.Mutex
	subscriptions map[string][]func([]byte)
}

func newGamificationTestPubSub() *gamificationTestPubSub {
	return &gamificationTestPubSub{
		subscriptions: make(map[string][]func([]byte)),
	}
}

func (p *gamificationTestPubSub) Subscribe(channel string, onMessage func([]byte)) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subscriptions[channel] = append(p.subscriptions[channel], onMessage)
}

func (p *gamificationTestPubSub) Publish(channel string, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	p.mu.Lock()
	subscribers := append([]func([]byte){}, p.subscriptions[channel]...)
	p.mu.Unlock()

	for _, subscriber := range subscribers {
		subscriber(data)
	}

	return nil
}

func setupGamificationRealtimeTest(t *testing.T) (*httptest.Server, *bun.DB) {
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
			oauth2_access_tokens,
			oauth2_refresh_tokens,
			oauth2_auth_codes,
			oauth2_oidc_sessions,
			users,
			oauth2_clients
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		t.Fatalf("failed to truncate realtime gamification tables: %v", err)
	}

	router := chi.NewRouter()
	authService := oauth2.NewOauth2SessionAuthServiceFromDB(testutils.TestTokenKey, db)
	pubSub := newGamificationTestPubSub()

	router.Route("/api/v1", func(r chi.Router) {
		wsModule := ws.NewWsModule(ws.WsModuleDeps{
			Router:         r,
			AuthService:    authService,
			TokenKey:       testutils.TestTokenKey,
			Server:         pubSub,
			OriginPatterns: []string{gamificationTestWSOrigin},
		})
		invitesModule := invites.NewInvitesModule(invites.InvitesModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
		})
		gamificationModule := gamification.NewGamificationModule(gamification.GamificationModuleDeps{
			Router:      r,
			DB:          db,
			AuthService: authService,
			WsService:   wsModule.Service,
		})
		rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:         r,
			DB:             db,
			WsService:      wsModule.Service,
			AuthService:    authService,
			InvitesService: invitesModule.Service,
			RewardService:  gamificationModule.Service,
		})
	})

	return httptest.NewServer(router), db
}

func connectGamificationWS(t *testing.T, serverURL, accessToken string) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(serverURL, "http") + "/api/v1/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Authorization": []string{"Bearer " + accessToken},
			"Origin":        []string{gamificationTestWSOrigin},
		},
	})
	if err != nil {
		t.Fatalf("failed to connect websocket: %v", err)
	}

	return conn
}

func readGamificationEvent(t *testing.T, conn *websocket.Conn, eventType string) ws.Event {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("failed to read websocket event %s: %v", eventType, err)
		}

		event := ws.Event{}
		if err := json.Unmarshal(data, &event); err != nil {
			t.Fatalf("failed to decode websocket event: %v", err)
		}

		if event.Type == eventType {
			return event
		}
	}
}

func TestFinishRoom_EmitsRewardEventToConnectedUser(t *testing.T) {
	server, db := setupGamificationRealtimeTest(t)
	defer server.Close()
	defer db.Close()

	adminToken, adminUserID := createGamificationAccessToken(t, db, "ws-admin@example.com")
	memberToken, memberUserID := createGamificationAccessToken(t, db, "ws-member@example.com")

	roomID := uuid.NewString()
	seedGamificationRoom(t, db, roomID, adminUserID, "ACTIVE", time.Now().UTC())

	adminParticipantID := seedGamificationParticipant(t, db, roomID, adminUserID, "ADMIN")
	memberParticipantID := seedGamificationParticipant(t, db, roomID, memberUserID, "MEMBER")
	finalEstimate := "5"
	taskID := seedGamificationTask(t, db, roomID, "ESTIMATED", false, &finalEstimate)
	seedGamificationTaskRound(t, db, taskID, 1, "REVEALED")
	seedGamificationVote(t, db, taskID, adminParticipantID, 1, "5")
	seedGamificationVote(t, db, taskID, memberParticipantID, 1, "5")

	memberConn := connectGamificationWS(t, server.URL, memberToken)
	defer memberConn.Close(websocket.StatusNormalClosure, "")
	readGamificationEvent(t, memberConn, ws.EventTypeHello)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, strings.NewReader(`{"status":"FINISHED"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.Config.Handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected room finish to succeed, got %d: %s", rr.Code, rr.Body.String())
	}

	event := readGamificationEvent(t, memberConn, gamification.SessionRewardedEvent)
	payload := struct {
		RoomID                    string `json:"roomId"`
		SessionsParticipatedDelta int    `json:"sessionsParticipatedDelta"`
		TasksEstimatedDelta       int    `json:"tasksEstimatedDelta"`
		XPGained                  int    `json:"xpGained"`
		PreviousXP                int    `json:"previousXp"`
		CurrentXP                 int    `json:"currentXp"`
	}{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("failed to decode reward payload: %v", err)
	}

	if payload.RoomID != roomID {
		t.Fatalf("expected reward room id %s, got %s", roomID, payload.RoomID)
	}
	if payload.SessionsParticipatedDelta != 1 || payload.TasksEstimatedDelta != 1 {
		t.Fatalf("unexpected reward deltas: %+v", payload)
	}
	if payload.XPGained != 13 || payload.PreviousXP != 0 || payload.CurrentXP != 13 {
		t.Fatalf("unexpected reward xp payload: %+v", payload)
	}
}
