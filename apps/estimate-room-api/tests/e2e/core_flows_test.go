package e2e

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/coder/websocket"
	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
)

type e2eUserMeResponse struct {
	ID    string  `json:"id"`
	Email *string `json:"email"`
}

type e2eTaskResponse struct {
	TaskID             string  `json:"TaskID"`
	Title              string  `json:"Title"`
	Status             string  `json:"Status"`
	IsActive           bool    `json:"IsActive"`
	FinalEstimateValue *string `json:"FinalEstimateValue"`
}

type e2eVotesRevealedPayload struct {
	TaskID      string `json:"taskId"`
	RoundNumber int    `json:"roundNumber"`
	AllVoted    bool   `json:"allVoted"`
	Summary     struct {
		TotalVotes int            `json:"totalVotes"`
		Counts     map[string]int `json:"counts"`
	} `json:"summary"`
}

func TestAuthFlow_LoginExchangeTokenAndReadProfile(t *testing.T) {
	app := setupE2EApp(t)

	accessToken := app.loginAndGetAccessToken(t, "alice@example.com", "password123")

	resp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/users/me",
		"",
		accessToken,
	)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when loading current user, got %d: %s", resp.StatusCode, readBody(t, resp))
	}

	me := decodeJSON[e2eUserMeResponse](t, resp)
	if me.ID == "" {
		t.Fatal("expected user id in /users/me response")
	}
	if me.Email == nil || *me.Email != "alice@example.com" {
		t.Fatalf("expected alice@example.com in /users/me response, got %#v", me.Email)
	}
}

func TestCreateRoomFlow(t *testing.T) {
	app := setupE2EApp(t)

	ownerToken := app.loginAndGetAccessToken(t, "owner@example.com", "password123")

	createRoomResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/rooms/",
		`{"name":"Flow Room","createShareLink":true}`,
		ownerToken,
	)
	if createRoomResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating room, got %d: %s", createRoomResp.StatusCode, readBody(t, createRoomResp))
	}

	createdRoom := decodeJSON[e2eCreateRoomResponse](t, createRoomResp)
	if createdRoom.Room.RoomID == "" {
		t.Fatal("expected room id in create room response")
	}
	if createdRoom.InviteToken == "" || createdRoom.ShareLink == nil {
		t.Fatalf("expected share link metadata in create room response, got %#v", createdRoom.ShareLink)
	}
	if len(createdRoom.Room.Participants) != 1 {
		t.Fatalf("expected exactly one room participant after creation, got %d", len(createdRoom.Room.Participants))
	}
	if createdRoom.Room.Participants[0].Role != "ADMIN" {
		t.Fatalf("expected creator to be ADMIN, got %s", createdRoom.Room.Participants[0].Role)
	}

	getRoomResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/rooms/"+createdRoom.Room.RoomID,
		"",
		ownerToken,
	)
	if getRoomResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when reading created room, got %d: %s", getRoomResp.StatusCode, readBody(t, getRoomResp))
	}

	room := decodeJSON[e2eRoomResponse](t, getRoomResp)
	if room.RoomID != createdRoom.Room.RoomID {
		t.Fatalf("expected room id %s, got %s", createdRoom.Room.RoomID, room.RoomID)
	}
}

func TestVotingRoundFlow(t *testing.T) {
	app := setupE2EApp(t)

	adminToken := app.loginAndGetAccessToken(t, "admin@example.com", "password123")
	memberToken := app.loginAndGetAccessToken(t, "member@example.com", "password123")

	createRoomResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/rooms/",
		`{"name":"Realtime Flow","inviteEmails":["member@example.com"]}`,
		adminToken,
	)
	if createRoomResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating realtime room, got %d: %s", createRoomResp.StatusCode, readBody(t, createRoomResp))
	}

	roomCreate := decodeJSON[e2eCreateRoomResponse](t, createRoomResp)
	if len(roomCreate.EmailInvites) != 1 {
		t.Fatalf("expected 1 room invite, got %d", len(roomCreate.EmailInvites))
	}

	acceptInviteResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/invites/"+roomCreate.EmailInvites[0].Token+"/accept",
		"",
		memberToken,
	)
	if acceptInviteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when accepting room invite, got %d: %s", acceptInviteResp.StatusCode, readBody(t, acceptInviteResp))
	}
	_ = decodeJSON[e2eRoomJoinResponse](t, acceptInviteResp)

	createTaskResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/rooms/"+roomCreate.Room.RoomID+"/tasks/",
		`{"title":"Estimate API latency"}`,
		adminToken,
	)
	if createTaskResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating task, got %d: %s", createTaskResp.StatusCode, readBody(t, createTaskResp))
	}

	task := decodeJSON[e2eTaskResponse](t, createTaskResp)
	if task.TaskID == "" {
		t.Fatal("expected task id in task create response")
	}

	adminConn := connectWS(t, app.server.URL, adminToken)
	defer adminConn.Close(websocket.StatusNormalClosure, "")

	memberConn := connectWS(t, app.server.URL, memberToken)
	defer memberConn.Close(websocket.StatusNormalClosure, "")

	joinRoom(t, adminConn, roomCreate.Room.RoomID)
	joinRoom(t, memberConn, roomCreate.Room.RoomID)

	writeEvent(t, adminConn, ws.Event{
		Type:   rooms.RoomsTaskSetCurrent,
		RoomID: roomCreate.Room.RoomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"taskId": task.TaskID,
		}),
	})

	_ = readUntilEvent(t, adminConn, rooms.RoomsTaskCurrentChanged)

	writeEvent(t, memberConn, ws.Event{
		Type:   rooms.RoomsVoteCast,
		RoomID: roomCreate.Room.RoomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"value": "5",
		}),
	})

	allCastEvent := readUntilEvent(t, adminConn, rooms.RoomsVotesAllCast)
	var allCastPayload struct {
		TaskID      string `json:"taskId"`
		RoundNumber int    `json:"roundNumber"`
	}
	if err := json.Unmarshal(allCastEvent.Payload, &allCastPayload); err != nil {
		t.Fatalf("failed to decode all-cast payload: %v", err)
	}
	if allCastPayload.TaskID != task.TaskID {
		t.Fatalf("expected all-cast task id %s, got %s", task.TaskID, allCastPayload.TaskID)
	}
	if allCastPayload.RoundNumber != 1 {
		t.Fatalf("expected round number 1, got %d", allCastPayload.RoundNumber)
	}

	writeEvent(t, adminConn, ws.Event{
		Type:   rooms.RoomsVoteReveal,
		RoomID: roomCreate.Room.RoomID,
	})

	revealedEvent := readUntilEvent(t, adminConn, rooms.RoomsVotesRevealed)
	revealedPayload := decodeEventPayload[e2eVotesRevealedPayload](t, revealedEvent.Payload)
	if !revealedPayload.AllVoted {
		t.Fatal("expected revealed payload to report AllVoted=true")
	}
	if revealedPayload.Summary.TotalVotes != 1 {
		t.Fatalf("expected 1 revealed vote, got %d", revealedPayload.Summary.TotalVotes)
	}
	if revealedPayload.Summary.Counts["5"] != 1 {
		t.Fatalf("expected one revealed vote for 5, got %#v", revealedPayload.Summary.Counts)
	}

	writeEvent(t, adminConn, ws.Event{
		Type:   rooms.RoomsTaskFinalize,
		RoomID: roomCreate.Room.RoomID,
		Payload: mustMarshalJSON(t, map[string]string{
			"value": "5",
		}),
	})

	finalizedEvent := readUntilEvent(t, adminConn, rooms.RoomsTaskFinalized)
	var finalizedPayload struct {
		TaskID             string `json:"taskId"`
		FinalEstimateValue string `json:"finalEstimateValue"`
		Status             string `json:"status"`
	}
	if err := json.Unmarshal(finalizedEvent.Payload, &finalizedPayload); err != nil {
		t.Fatalf("failed to decode finalized payload: %v", err)
	}
	if finalizedPayload.TaskID != task.TaskID {
		t.Fatalf("expected finalized task id %s, got %s", task.TaskID, finalizedPayload.TaskID)
	}
	if finalizedPayload.FinalEstimateValue != "5" || finalizedPayload.Status != "ESTIMATED" {
		t.Fatalf("unexpected finalized payload: %#v", finalizedPayload)
	}

	getTaskResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/rooms/"+roomCreate.Room.RoomID+"/tasks/"+task.TaskID,
		"",
		adminToken,
	)
	if getTaskResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when reading finalized task, got %d: %s", getTaskResp.StatusCode, readBody(t, getTaskResp))
	}

	finalTask := decodeJSON[e2eTaskResponse](t, getTaskResp)
	if finalTask.FinalEstimateValue == nil || *finalTask.FinalEstimateValue != "5" {
		t.Fatalf("expected persisted final estimate 5, got %#v", finalTask.FinalEstimateValue)
	}
	if finalTask.Status != "ESTIMATED" || finalTask.IsActive {
		t.Fatalf("expected ESTIMATED inactive task, got status=%s active=%v", finalTask.Status, finalTask.IsActive)
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

func decodeEventPayload[T any](t *testing.T, payload json.RawMessage) T {
	t.Helper()

	var value T
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("failed to decode websocket payload: %v", err)
	}

	return value
}
