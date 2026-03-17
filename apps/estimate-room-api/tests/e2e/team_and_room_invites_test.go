package e2e

import (
	"net/http"
	"testing"
)

type e2eInvitationResponse struct {
	InvitationID  string  `json:"invitationId"`
	Kind          string  `json:"kind"`
	Status        string  `json:"status"`
	TeamID        *string `json:"teamId"`
	RoomID        *string `json:"roomId"`
	InvitedEmail  *string `json:"invitedEmail"`
	Token         string  `json:"token"`
	InvitedUserID *string `json:"invitedUserId"`
}

type e2eTeamMemberResponse struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	User   struct {
		UserID string  `json:"userId"`
		Email  *string `json:"email"`
	} `json:"user"`
}

type e2eTeamDetailResponse struct {
	TeamID  string                  `json:"teamId"`
	Name    string                  `json:"name"`
	Members []e2eTeamMemberResponse `json:"members"`
}

type e2eRoomResponse struct {
	RoomID       string `json:"RoomID"`
	Name         string `json:"Name"`
	Participants []struct {
		Role      string  `json:"Role"`
		UserID    *string `json:"UserID"`
		GuestName *string `json:"GuestName"`
	} `json:"Participants"`
}

type e2eRoomJoinResponse struct {
	Room        e2eRoomResponse `json:"room"`
	Participant struct {
		Role      string  `json:"Role"`
		UserID    *string `json:"UserID"`
		GuestName *string `json:"GuestName"`
	} `json:"participant"`
}

type e2eCreateRoomResponse struct {
	Room              e2eRoomResponse         `json:"room"`
	EmailInvites      []e2eInvitationResponse `json:"emailInvites"`
	ShareLink         *e2eInvitationResponse  `json:"shareLink"`
	InviteToken       string                  `json:"inviteToken"`
	SkippedRecipients []map[string]any        `json:"skippedRecipients"`
}

func TestTeamInviteFlow(t *testing.T) {
	app := setupE2EApp(t)

	ownerToken := app.loginAndGetAccessToken(t, "owner@example.com", "password123")
	memberToken := app.loginAndGetAccessToken(t, "member@example.com", "password123")

	createTeamResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/teams/",
		`{"name":"Platform Team"}`,
		ownerToken,
	)
	if createTeamResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating team, got %d: %s", createTeamResp.StatusCode, readBody(t, createTeamResp))
	}

	team := decodeJSON[e2eTeamDetailResponse](t, createTeamResp)
	if team.TeamID == "" {
		t.Fatal("expected team id in create response")
	}

	createInviteResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/teams/"+team.TeamID+"/invites",
		`{"emails":["member@example.com"]}`,
		ownerToken,
	)
	if createInviteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating team invite, got %d: %s", createInviteResp.StatusCode, readBody(t, createInviteResp))
	}

	invites := decodeJSON[[]e2eInvitationResponse](t, createInviteResp)
	if len(invites) != 1 {
		t.Fatalf("expected 1 team invite, got %d", len(invites))
	}
	if invites[0].Kind != "TEAM_MEMBER" {
		t.Fatalf("expected TEAM_MEMBER invite, got %s", invites[0].Kind)
	}
	if invites[0].Token == "" {
		t.Fatal("expected token in team invite response")
	}

	previewResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/invites/"+invites[0].Token,
		"",
		"",
	)
	if previewResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when previewing team invite, got %d: %s", previewResp.StatusCode, readBody(t, previewResp))
	}

	preview := decodeJSON[e2eInvitationResponse](t, previewResp)
	if preview.Kind != "TEAM_MEMBER" {
		t.Fatalf("expected TEAM_MEMBER preview, got %s", preview.Kind)
	}
	if preview.TeamID == nil || *preview.TeamID != team.TeamID {
		t.Fatalf("expected preview team id %s, got %#v", team.TeamID, preview.TeamID)
	}

	acceptResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/invites/"+invites[0].Token+"/accept",
		"",
		memberToken,
	)
	if acceptResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when accepting team invite, got %d: %s", acceptResp.StatusCode, readBody(t, acceptResp))
	}

	accepted := decodeJSON[e2eInvitationResponse](t, acceptResp)
	if accepted.Status != "ACCEPTED" {
		t.Fatalf("expected ACCEPTED team invite, got %s", accepted.Status)
	}

	getTeamResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/teams/"+team.TeamID,
		"",
		memberToken,
	)
	if getTeamResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when member reads team, got %d: %s", getTeamResp.StatusCode, readBody(t, getTeamResp))
	}

	updatedTeam := decodeJSON[e2eTeamDetailResponse](t, getTeamResp)
	if len(updatedTeam.Members) != 2 {
		t.Fatalf("expected 2 team members after accept, got %d", len(updatedTeam.Members))
	}
}

func TestRoomInviteFlow(t *testing.T) {
	app := setupE2EApp(t)

	ownerToken := app.loginAndGetAccessToken(t, "owner@example.com", "password123")
	memberToken := app.loginAndGetAccessToken(t, "member@example.com", "password123")
	outsideToken := app.loginAndGetAccessToken(t, "outside@example.com", "password123")

	createTeamResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/teams/",
		`{"name":"Delivery Team"}`,
		ownerToken,
	)
	if createTeamResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating team, got %d: %s", createTeamResp.StatusCode, readBody(t, createTeamResp))
	}

	team := decodeJSON[e2eTeamDetailResponse](t, createTeamResp)

	createTeamInviteResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/teams/"+team.TeamID+"/invites",
		`{"emails":["member@example.com"]}`,
		ownerToken,
	)
	if createTeamInviteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating member invite, got %d: %s", createTeamInviteResp.StatusCode, readBody(t, createTeamInviteResp))
	}

	teamInvites := decodeJSON[[]e2eInvitationResponse](t, createTeamInviteResp)
	if len(teamInvites) != 1 {
		t.Fatalf("expected 1 team invite, got %d", len(teamInvites))
	}

	acceptTeamInviteResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/invites/"+teamInvites[0].Token+"/accept",
		"",
		memberToken,
	)
	if acceptTeamInviteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when accepting team invite, got %d: %s", acceptTeamInviteResp.StatusCode, readBody(t, acceptTeamInviteResp))
	}
	_ = decodeJSON[e2eInvitationResponse](t, acceptTeamInviteResp)

	createRoomResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/rooms/",
		`{
			"name":"Sprint Planning",
			"inviteTeamId":"`+team.TeamID+`",
			"inviteEmails":["outside@example.com"],
			"createShareLink":true
		}`,
		ownerToken,
	)
	if createRoomResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when creating room, got %d: %s", createRoomResp.StatusCode, readBody(t, createRoomResp))
	}

	roomCreate := decodeJSON[e2eCreateRoomResponse](t, createRoomResp)
	if roomCreate.Room.RoomID == "" {
		t.Fatal("expected room id in create room response")
	}
	if roomCreate.ShareLink == nil || roomCreate.InviteToken == "" {
		t.Fatalf("expected share link in create room response, got %#v", roomCreate.ShareLink)
	}
	if len(roomCreate.EmailInvites) != 2 {
		t.Fatalf("expected 2 room email invites, got %d", len(roomCreate.EmailInvites))
	}

	var outsideInviteToken string
	for _, invite := range roomCreate.EmailInvites {
		if invite.Kind != "ROOM_EMAIL" {
			t.Fatalf("expected ROOM_EMAIL invite, got %s", invite.Kind)
		}
		if invite.RoomID == nil || *invite.RoomID != roomCreate.Room.RoomID {
			t.Fatalf("expected invite room id %s, got %#v", roomCreate.Room.RoomID, invite.RoomID)
		}
		if invite.InvitedEmail != nil && *invite.InvitedEmail == "outside@example.com" {
			outsideInviteToken = invite.Token
		}
	}
	if outsideInviteToken == "" {
		t.Fatalf("expected room invite for outside@example.com, got %#v", roomCreate.EmailInvites)
	}

	previewResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/invites/"+outsideInviteToken,
		"",
		"",
	)
	if previewResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when previewing room invite, got %d: %s", previewResp.StatusCode, readBody(t, previewResp))
	}

	preview := decodeJSON[e2eInvitationResponse](t, previewResp)
	if preview.Kind != "ROOM_EMAIL" {
		t.Fatalf("expected ROOM_EMAIL preview, got %s", preview.Kind)
	}
	if preview.RoomID == nil || *preview.RoomID != roomCreate.Room.RoomID {
		t.Fatalf("expected preview room id %s, got %#v", roomCreate.Room.RoomID, preview.RoomID)
	}

	acceptRoomInviteResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/invites/"+outsideInviteToken+"/accept",
		"",
		outsideToken,
	)
	if acceptRoomInviteResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when accepting room email invite, got %d: %s", acceptRoomInviteResp.StatusCode, readBody(t, acceptRoomInviteResp))
	}

	joinedMember := decodeJSON[e2eRoomJoinResponse](t, acceptRoomInviteResp)
	if joinedMember.Participant.Role != "MEMBER" {
		t.Fatalf("expected MEMBER role after room email accept, got %s", joinedMember.Participant.Role)
	}

	getRoomAsMemberResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodGet,
		app.server.URL+"/api/v1/rooms/"+roomCreate.Room.RoomID,
		"",
		outsideToken,
	)
	if getRoomAsMemberResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when invited member reads room, got %d: %s", getRoomAsMemberResp.StatusCode, readBody(t, getRoomAsMemberResp))
	}
	_ = decodeJSON[e2eRoomResponse](t, getRoomAsMemberResp)

	guestJoinResp := doJSONRequest(
		t,
		app.server.Client(),
		http.MethodPost,
		app.server.URL+"/api/v1/invites/"+roomCreate.InviteToken+"/accept",
		`{"guestName":"Guest One"}`,
		"",
	)
	if guestJoinResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when guest joins via share link, got %d: %s", guestJoinResp.StatusCode, readBody(t, guestJoinResp))
	}

	guestJoin := decodeJSON[e2eRoomJoinResponse](t, guestJoinResp)
	if guestJoin.Participant.Role != "GUEST" {
		t.Fatalf("expected GUEST role after share-link join, got %s", guestJoin.Participant.Role)
	}
	if guestJoin.Participant.GuestName == nil || *guestJoin.Participant.GuestName != "Guest One" {
		t.Fatalf("expected guest name Guest One, got %#v", guestJoin.Participant.GuestName)
	}

	guestReadReq, err := http.NewRequest(http.MethodGet, app.server.URL+"/api/v1/rooms/"+roomCreate.Room.RoomID, nil)
	if err != nil {
		t.Fatalf("failed to build guest read request: %v", err)
	}
	for _, cookie := range guestJoinResp.Cookies() {
		guestReadReq.AddCookie(cookie)
	}

	guestReadResp, err := app.server.Client().Do(guestReadReq)
	if err != nil {
		t.Fatalf("failed to read room as guest: %v", err)
	}
	if guestReadResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 when guest reads room, got %d: %s", guestReadResp.StatusCode, readBody(t, guestReadResp))
	}
	_ = decodeJSON[e2eRoomResponse](t, guestReadResp)
}
