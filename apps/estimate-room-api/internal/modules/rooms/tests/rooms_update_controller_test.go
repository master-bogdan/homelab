package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/master-bogdan/estimate-room-api/internal/modules/gamification"
	roomsdto "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/dto"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type failingRewardService struct {
	applyErr error
}

func (s *failingRewardService) ApplyRoomTerminalRewards(
	ctx context.Context,
	db bun.IDB,
	room *roomsmodels.RoomsModel,
) ([]gamification.AppliedRoomReward, error) {
	return nil, s.applyErr
}

func (s *failingRewardService) NotifyAppliedRewards(
	ctx context.Context,
	rewards []gamification.AppliedRoomReward,
) error {
	return nil
}

func TestUpdateRoom_ChangesFieldsForCreator(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"name":"Renamed room"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var room roomsdto.RoomResponse
	if err := json.NewDecoder(rr.Body).Decode(&room); err != nil {
		t.Fatalf("failed to decode room response: %v", err)
	}

	if room.Name != "Renamed room" {
		t.Fatalf("expected updated name, got %s", room.Name)
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

	var room roomsdto.RoomResponse
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

func TestUpdateRoom_CannotRenameFinishedRoom(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET status = 'FINISHED', finished_at = NOW()
		WHERE room_id = $1
	`, roomID); err != nil {
		t.Fatalf("failed to finish room: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"name":"Renamed after finish"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}

	room := loadRoomForUpdateTest(t, db, roomID)
	if room.Name != "Test Room" {
		t.Fatalf("expected finished room name to stay unchanged, got %s", room.Name)
	}
	if room.Status != "FINISHED" {
		t.Fatalf("expected room to stay FINISHED, got %s", room.Status)
	}
}

func TestUpdateRoom_CannotReopenFinishedRoom(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET status = 'FINISHED', finished_at = NOW()
		WHERE room_id = $1
	`, roomID); err != nil {
		t.Fatalf("failed to finish room: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"status":"ACTIVE"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}

	room := loadRoomForUpdateTest(t, db, roomID)
	if room.Status != "FINISHED" {
		t.Fatalf("expected room to stay FINISHED, got %s", room.Status)
	}
	if room.FinishedAt == nil {
		t.Fatal("expected finishedAt to stay set")
	}
}

func TestUpdateRoom_NoopPatchOnFinishedRoomReturnsBadRequest(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET status = 'FINISHED', finished_at = NOW()
		WHERE room_id = $1
	`, roomID); err != nil {
		t.Fatalf("failed to finish room: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"status":"FINISHED"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestUpdateRoom_CannotEditExpiredRoom(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET status = 'EXPIRED', finished_at = NOW()
		WHERE room_id = $1
	`, roomID); err != nil {
		t.Fatalf("failed to expire room: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"name":"Expired rename"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}

	room := loadRoomForUpdateTest(t, db, roomID)
	if room.Name != "Test Room" {
		t.Fatalf("expected expired room name to stay unchanged, got %s", room.Name)
	}
	if room.Status != "EXPIRED" {
		t.Fatalf("expected room to stay EXPIRED, got %s", room.Status)
	}
}

func TestUpdateRoom_CannotSetExpiredStatusManually(t *testing.T) {
	router, db := setupRoomsTasksTest(t)
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"status":"EXPIRED"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d: %s", rr.Code, rr.Body.String())
	}

	room := loadRoomForUpdateTest(t, db, roomID)
	if room.Status != "ACTIVE" {
		t.Fatalf("expected room to stay ACTIVE, got %s", room.Status)
	}
	if room.FinishedAt != nil {
		t.Fatalf("expected finishedAt to remain nil, got %v", room.FinishedAt)
	}
}

func TestUpdateRoom_FinishPersistsWhenRewardApplicationFails(t *testing.T) {
	router, db := setupRoomsTasksTestWithRewardService(t, &failingRewardService{
		applyErr: context.DeadlineExceeded,
	})
	defer db.Close()

	accessToken, userID := createAccessToken(t, db)
	roomID := seedRoom(t, db, userID)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/rooms/"+roomID, bytes.NewReader([]byte(`{"status":"FINISHED"}`)))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when reward application fails, got %d: %s", rr.Code, rr.Body.String())
	}

	room := loadRoomForUpdateTest(t, db, roomID)
	if room.Status != "FINISHED" {
		t.Fatalf("expected room to persist FINISHED status, got %s", room.Status)
	}
	if room.FinishedAt == nil {
		t.Fatal("expected finishedAt to stay set even when rewards fail")
	}
}

func loadRoomForUpdateTest(t *testing.T, db *bun.DB, roomID string) *roomsmodels.RoomsModel {
	t.Helper()

	var room roomsmodels.RoomsModel
	err := db.NewSelect().
		Model(&room).
		Where("room_id = ?", roomID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		t.Fatalf("failed to load room %s: %v", roomID, err)
	}

	return &room
}
