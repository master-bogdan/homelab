package tests

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	testutils "github.com/master-bogdan/estimate-room-api/internal/pkg/test"
	"github.com/uptrace/bun"
)

type testPubSub struct {
	mu            sync.Mutex
	subscriptions map[string][]func([]byte)
}

func newTestPubSub() *testPubSub {
	return &testPubSub{
		subscriptions: make(map[string][]func([]byte)),
	}
}

func (p *testPubSub) Subscribe(channel string, onMessage func([]byte)) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.subscriptions[channel] = append(p.subscriptions[channel], onMessage)
}

func (p *testPubSub) Publish(channel string, message any) error {
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

func setupRoomsRealtimeTest(t *testing.T) (*httptest.Server, *bun.DB) {
	t.Helper()

	db := testutils.SetupTestDB(t)

	_, err := db.ExecContext(context.Background(), `
		TRUNCATE TABLE
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
		t.Fatalf("failed to truncate tables: %v", err)
	}

	router := chi.NewRouter()
	authService := oauth2.NewAuthServiceFromDB(testutils.TestTokenKey, db)
	pubSub := newTestPubSub()

	router.Route("/api/v1", func(r chi.Router) {
		wsModule := ws.NewWsModule(ws.WsModuleDeps{
			Router:      r,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
			Server:      pubSub,
		})

		rooms.NewRoomsModule(rooms.RoomsModuleDeps{
			Router:      r,
			DB:          db,
			WsService:   wsModule.Service,
			AuthService: authService,
			TokenKey:    testutils.TestTokenKey,
		})
	})

	return httptest.NewServer(router), db
}

func seedMemberParticipant(t *testing.T, db *bun.DB, roomID, userID string) string {
	t.Helper()

	participantID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO room_participants (room_participants_id, room_id, user_id, role)
		VALUES ($1, $2, $3, 'MEMBER')
	`, participantID, roomID, userID)
	if err != nil {
		t.Fatalf("failed to insert member participant: %v", err)
	}

	return participantID
}

func seedTask(t *testing.T, db *bun.DB, roomID, title string) string {
	t.Helper()

	taskID := uuid.NewString()
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO tasks (task_id, room_id, title, status, is_active)
		VALUES ($1, $2, $3, 'PENDING', FALSE)
	`, taskID, roomID, title)
	if err != nil {
		t.Fatalf("failed to insert task: %v", err)
	}

	return taskID
}

func connectWS(t *testing.T, serverURL, accessToken string) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(serverURL, "http") + "/api/v1/ws?token=" + url.QueryEscape(accessToken)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect websocket: %v", err)
	}

	return conn
}

func readUntilEvent(t *testing.T, conn *websocket.Conn, eventType string) ws.Event {
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

func writeEvent(t *testing.T, conn *websocket.Conn, event ws.Event) {
	t.Helper()

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal websocket event: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		t.Fatalf("failed to write websocket event: %v", err)
	}
}

func joinRoom(t *testing.T, conn *websocket.Conn, roomID string) {
	t.Helper()

	readUntilEvent(t, conn, ws.EventTypeHello)
	writeEvent(t, conn, ws.Event{
		Type:   rooms.RoomsJoin,
		RoomID: roomID,
	})
	readUntilEvent(t, conn, rooms.RoomsSnapshot)
}

func decodePayload[T any](t *testing.T, payload json.RawMessage) T {
	t.Helper()

	var value T
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	return value
}

func TestRoomsVoting_AllEligibleVotesEmitAllCastEvent(t *testing.T) {
	server, db := setupRoomsRealtimeTest(t)
	defer server.Close()
	defer db.Close()

	adminToken, adminUserID := createAccessToken(t, db)
	roomID := seedRoom(t, db, adminUserID)
	taskID := seedTask(t, db, roomID, "Realtime task")

	memberTokenOne, memberUserIDOne := createAccessToken(t, db)
	memberTokenTwo, memberUserIDTwo := createAccessToken(t, db)
	memberParticipantIDOne := seedMemberParticipant(t, db, roomID, memberUserIDOne)
	memberParticipantIDTwo := seedMemberParticipant(t, db, roomID, memberUserIDTwo)

	adminConn := connectWS(t, server.URL, adminToken)
	defer adminConn.Close(websocket.StatusNormalClosure, "")
	memberConnOne := connectWS(t, server.URL, memberTokenOne)
	defer memberConnOne.Close(websocket.StatusNormalClosure, "")
	memberConnTwo := connectWS(t, server.URL, memberTokenTwo)
	defer memberConnTwo.Close(websocket.StatusNormalClosure, "")

	joinRoom(t, adminConn, roomID)
	joinRoom(t, memberConnOne, roomID)
	joinRoom(t, memberConnTwo, roomID)

	writeEvent(t, adminConn, ws.Event{
		Type:   rooms.RoomsTaskSetCurrent,
		RoomID: roomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"taskId": taskID,
		}),
	})

	currentChangedEvent := readUntilEvent(t, adminConn, rooms.RoomsTaskCurrentChanged)
	currentChangedPayload := decodePayload[struct {
		CurrentTaskID          string   `json:"currentTaskId"`
		RoundNumber            int      `json:"roundNumber"`
		EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
	}](t, currentChangedEvent.Payload)

	if currentChangedPayload.CurrentTaskID != taskID {
		t.Fatalf("expected current task %s, got %s", taskID, currentChangedPayload.CurrentTaskID)
	}
	if !sameStringSet(currentChangedPayload.EligibleParticipantIDs, []string{memberParticipantIDOne, memberParticipantIDTwo}) {
		t.Fatalf("expected eligible participants %v, got %v", []string{memberParticipantIDOne, memberParticipantIDTwo}, currentChangedPayload.EligibleParticipantIDs)
	}

	writeEvent(t, memberConnOne, ws.Event{
		Type:   rooms.RoomsVoteCast,
		RoomID: roomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"value": "3",
		}),
	})
	readUntilEvent(t, adminConn, rooms.RoomsVoteStatusChanged)

	writeEvent(t, memberConnTwo, ws.Event{
		Type:   rooms.RoomsVoteCast,
		RoomID: roomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"value": "5",
		}),
	})

	allCastEvent := readUntilEvent(t, adminConn, rooms.RoomsVotesAllCast)
	allCastPayload := decodePayload[struct {
		TaskID                 string   `json:"taskId"`
		RoundNumber            int      `json:"roundNumber"`
		EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
		VotedParticipantIDs    []string `json:"votedParticipantIds"`
	}](t, allCastEvent.Payload)

	if allCastPayload.TaskID != taskID {
		t.Fatalf("expected all-cast task %s, got %s", taskID, allCastPayload.TaskID)
	}
	if allCastPayload.RoundNumber != 1 {
		t.Fatalf("expected round number 1, got %d", allCastPayload.RoundNumber)
	}
	if !sameStringSet(allCastPayload.EligibleParticipantIDs, []string{memberParticipantIDOne, memberParticipantIDTwo}) {
		t.Fatalf("expected eligible participants %v, got %v", []string{memberParticipantIDOne, memberParticipantIDTwo}, allCastPayload.EligibleParticipantIDs)
	}
	if !sameStringSet(allCastPayload.VotedParticipantIDs, []string{memberParticipantIDOne, memberParticipantIDTwo}) {
		t.Fatalf("expected voted participants %v, got %v", []string{memberParticipantIDOne, memberParticipantIDTwo}, allCastPayload.VotedParticipantIDs)
	}
}

func TestRoomsVoting_LateJoinerDoesNotEnterCurrentRoundEligibility(t *testing.T) {
	server, db := setupRoomsRealtimeTest(t)
	defer server.Close()
	defer db.Close()

	adminToken, adminUserID := createAccessToken(t, db)
	roomID := seedRoom(t, db, adminUserID)
	taskID := seedTask(t, db, roomID, "Realtime task")

	memberTokenOne, memberUserIDOne := createAccessToken(t, db)
	memberTokenTwo, memberUserIDTwo := createAccessToken(t, db)
	memberParticipantIDOne := seedMemberParticipant(t, db, roomID, memberUserIDOne)
	memberParticipantIDTwo := seedMemberParticipant(t, db, roomID, memberUserIDTwo)

	adminConn := connectWS(t, server.URL, adminToken)
	defer adminConn.Close(websocket.StatusNormalClosure, "")
	memberConnOne := connectWS(t, server.URL, memberTokenOne)
	defer memberConnOne.Close(websocket.StatusNormalClosure, "")

	joinRoom(t, adminConn, roomID)
	joinRoom(t, memberConnOne, roomID)

	writeEvent(t, adminConn, ws.Event{
		Type:   rooms.RoomsTaskSetCurrent,
		RoomID: roomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"taskId": taskID,
		}),
	})

	currentChangedEvent := readUntilEvent(t, adminConn, rooms.RoomsTaskCurrentChanged)
	currentChangedPayload := decodePayload[struct {
		EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
	}](t, currentChangedEvent.Payload)
	if !sameStringSet(currentChangedPayload.EligibleParticipantIDs, []string{memberParticipantIDOne}) {
		t.Fatalf("expected only the already connected member to be eligible, got %v", currentChangedPayload.EligibleParticipantIDs)
	}

	memberConnTwo := connectWS(t, server.URL, memberTokenTwo)
	defer memberConnTwo.Close(websocket.StatusNormalClosure, "")
	joinRoom(t, memberConnTwo, roomID)

	writeEvent(t, memberConnTwo, ws.Event{
		Type:   rooms.RoomsVoteCast,
		RoomID: roomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"value": "5",
		}),
	})

	time.Sleep(200 * time.Millisecond)

	var lateJoinerVoteCount int
	err := db.QueryRowContext(context.Background(), `
		SELECT COUNT(*)
		FROM votes
		WHERE task_id = $1 AND participant_id = $2
	`, taskID, memberParticipantIDTwo).Scan(&lateJoinerVoteCount)
	if err != nil {
		t.Fatalf("failed to count late joiner votes: %v", err)
	}
	if lateJoinerVoteCount != 0 {
		t.Fatalf("expected no votes from late joiner, got %d", lateJoinerVoteCount)
	}

	writeEvent(t, memberConnOne, ws.Event{
		Type:   rooms.RoomsVoteCast,
		RoomID: roomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"value": "3",
		}),
	})

	allCastEvent := readUntilEvent(t, adminConn, rooms.RoomsVotesAllCast)
	allCastPayload := decodePayload[struct {
		EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
		VotedParticipantIDs    []string `json:"votedParticipantIds"`
	}](t, allCastEvent.Payload)

	if !sameStringSet(allCastPayload.EligibleParticipantIDs, []string{memberParticipantIDOne}) {
		t.Fatalf("expected late joiner to stay out of eligibility, got %v", allCastPayload.EligibleParticipantIDs)
	}
	if !sameStringSet(allCastPayload.VotedParticipantIDs, []string{memberParticipantIDOne}) {
		t.Fatalf("expected only the original eligible voter in voted list, got %v", allCastPayload.VotedParticipantIDs)
	}
}

func mustMarshalJSON(t *testing.T, value any) json.RawMessage {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	return data
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	normalizedLeft := append([]string(nil), left...)
	normalizedRight := append([]string(nil), right...)
	sort.Strings(normalizedLeft)
	sort.Strings(normalizedRight)

	for i := range normalizedLeft {
		if normalizedLeft[i] != normalizedRight[i] {
			return false
		}
	}

	return true
}
