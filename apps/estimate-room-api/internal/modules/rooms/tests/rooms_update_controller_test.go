package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

func TestUpdateRoom_ChangesFieldsForCreator(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"name":"Renamed room","allowGuests":true,"allowSpectators":true,"roundTimerSeconds":300}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var room roomsmodels.RoomsModel
	if err := json.NewDecoder(rr.Body).Decode(&room); err != nil {
		t.Fatalf("failed to decode room response: %v", err)
	}

	if room.Name != "Renamed room" {
		t.Fatalf("expected updated name, got %s", room.Name)
	}
	if !room.AllowGuests {
		t.Fatal("expected allowGuests to be true")
	}
	if !room.AllowSpectators {
		t.Fatal("expected allowSpectators to be true")
	}
	if room.RoundTimerSeconds != 300 {
		t.Fatalf("expected roundTimerSeconds 300, got %d", room.RoundTimerSeconds)
	}
}

func TestUpdateRoom_FinishingSetsFinishedAt(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"status":"FINISHED"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var room roomsmodels.RoomsModel
	if err := json.NewDecoder(rr.Body).Decode(&room); err != nil {
		t.Fatalf("failed to decode room response: %v", err)
	}

	if room.Status != "FINISHED" {
		t.Fatalf("expected FINISHED status, got %s", room.Status)
	}
	if room.FinishedAt == nil {
		t.Fatal("expected finishedAt to be set")
	}
}

func TestUpdateRoom_OnlyCreatorCanUpdate(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	_, ownerUserID := createAccessToken(t, db)
	roomID := seedRoom(t, db, ownerUserID)

	otherToken, _ := createAccessToken(t, db)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"status":"FINISHED"}`)))
	req.Header.Set("Authorization", "Bearer "+otherToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d: %s", rr.Code, rr.Body.String())
	}

	var httpErr apperrors.HttpError
	if err := json.NewDecoder(rr.Body).Decode(&httpErr); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if httpErr.Status != http.StatusForbidden {
		t.Fatalf("expected error status 403, got %d", httpErr.Status)
	}
}
