package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/modules/rooms"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
)

func TestRoomsExpiry_ExpireInactiveRoomsMarksStaleRoomsExpiredAndBroadcastsEvent(t *testing.T) {
	_, db := setupRoomsTasksTest(t)
	defer db.Close()

	_, staleAdminUserID := createAccessToken(t, db)
	staleRoomID := seedRoom(t, db, staleAdminUserID)
	_, freshAdminUserID := createAccessToken(t, db)
	freshRoomID := seedRoom(t, db, freshAdminUserID)

	staleActivityAt := time.Now().Add(-31 * time.Minute).UTC()
	freshActivityAt := time.Now().Add(-5 * time.Minute).UTC()

	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET last_activity_at = $2
		WHERE room_id = $1
	`, staleRoomID, staleActivityAt); err != nil {
		t.Fatalf("failed to set stale room activity time: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET last_activity_at = $2
		WHERE room_id = $1
	`, freshRoomID, freshActivityAt); err != nil {
		t.Fatalf("failed to set fresh room activity time: %v", err)
	}

	pubSub := newTestPubSub()
	wsService := ws.NewService(pubSub, "test-room-events")
	expiryService := rooms.NewRoomsExpiryService(roomsrepositories.NewRoomsRepository(db), wsService)

	events := make(chan ws.Event, 1)
	pubSub.Subscribe("test-room-events", func(data []byte) {
		event := ws.Event{}
		if err := json.Unmarshal(data, &event); err != nil {
			return
		}
		if event.Type == rooms.RoomsExpired {
			select {
			case events <- event:
			default:
			}
		}
	})

	expiredRooms, err := expiryService.ExpireInactiveRooms(time.Now().Add(-30 * time.Minute))
	if err != nil {
		t.Fatalf("failed to expire inactive rooms: %v", err)
	}
	if len(expiredRooms) != 1 {
		t.Fatalf("expected 1 expired room, got %d", len(expiredRooms))
	}
	if expiredRooms[0].RoomID != staleRoomID {
		t.Fatalf("expected stale room %s to expire, got %s", staleRoomID, expiredRooms[0].RoomID)
	}

	staleRoomRepo := roomsrepositories.NewRoomsRepository(db)
	staleRoom, err := staleRoomRepo.FindByID(staleRoomID)
	if err != nil {
		t.Fatalf("failed to load expired room: %v", err)
	}
	if staleRoom.Status != "EXPIRED" {
		t.Fatalf("expected stale room status EXPIRED, got %s", staleRoom.Status)
	}
	if staleRoom.FinishedAt == nil {
		t.Fatal("expected expired room to set finishedAt")
	}

	freshRoom, err := staleRoomRepo.FindByID(freshRoomID)
	if err != nil {
		t.Fatalf("failed to load fresh room: %v", err)
	}
	if freshRoom.Status != "ACTIVE" {
		t.Fatalf("expected fresh room to remain ACTIVE, got %s", freshRoom.Status)
	}

	select {
	case event := <-events:
		if event.RoomID != staleRoomID {
			t.Fatalf("expected room expired event for room %s, got %s", staleRoomID, event.RoomID)
		}

		payload := struct {
			RoomID string `json:"roomId"`
			Status string `json:"status"`
		}{}
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatalf("failed to decode room expired payload: %v", err)
		}
		if payload.RoomID != staleRoomID {
			t.Fatalf("expected payload room id %s, got %s", staleRoomID, payload.RoomID)
		}
		if payload.Status != "EXPIRED" {
			t.Fatalf("expected payload status EXPIRED, got %s", payload.Status)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("expected room expired realtime event")
	}
}

func TestRoomsExpiry_TouchActivityUpdatesOnlyActiveRooms(t *testing.T) {
	_, db := setupRoomsTasksTest(t)
	defer db.Close()

	_, activeAdminUserID := createAccessToken(t, db)
	activeRoomID := seedRoom(t, db, activeAdminUserID)
	_, expiredAdminUserID := createAccessToken(t, db)
	expiredRoomID := seedRoom(t, db, expiredAdminUserID)

	activeOldActivityAt := time.Now().Add(-2 * time.Hour).UTC()
	expiredOldActivityAt := time.Now().Add(-90 * time.Minute).UTC()

	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET last_activity_at = $2
		WHERE room_id = $1
	`, activeRoomID, activeOldActivityAt); err != nil {
		t.Fatalf("failed to set active room activity time: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET status = 'EXPIRED', finished_at = NOW(), last_activity_at = $2
		WHERE room_id = $1
	`, expiredRoomID, expiredOldActivityAt); err != nil {
		t.Fatalf("failed to set expired room activity time: %v", err)
	}

	expiryService := rooms.NewRoomsExpiryService(roomsrepositories.NewRoomsRepository(db), nil)
	expiryService.TouchActivity(activeRoomID)
	expiryService.TouchActivity(expiredRoomID)

	roomsRepo := roomsrepositories.NewRoomsRepository(db)

	activeRoom, err := roomsRepo.FindByID(activeRoomID)
	if err != nil {
		t.Fatalf("failed to load active room: %v", err)
	}
	if !activeRoom.LastActivityAt.After(activeOldActivityAt) {
		t.Fatalf("expected active room lastActivityAt to advance beyond %s, got %s", activeOldActivityAt, activeRoom.LastActivityAt)
	}

	expiredRoom, err := roomsRepo.FindByID(expiredRoomID)
	if err != nil {
		t.Fatalf("failed to load expired room: %v", err)
	}
	if expiredRoom.LastActivityAt.Sub(expiredOldActivityAt) > time.Second || expiredOldActivityAt.Sub(expiredRoom.LastActivityAt) > time.Second {
		t.Fatalf("expected expired room lastActivityAt to stay %s, got %s", expiredOldActivityAt, expiredRoom.LastActivityAt)
	}
	if expiredRoom.Status != "EXPIRED" {
		t.Fatalf("expected expired room to stay EXPIRED, got %s", expiredRoom.Status)
	}
}

func TestRoomsExpiry_IgnoreFinishedRooms(t *testing.T) {
	_, db := setupRoomsTasksTest(t)
	defer db.Close()

	_, finishedAdminUserID := createAccessToken(t, db)
	finishedRoomID := seedRoom(t, db, finishedAdminUserID)
	_, staleActiveAdminUserID := createAccessToken(t, db)
	staleActiveRoomID := seedRoom(t, db, staleActiveAdminUserID)

	oldActivityAt := time.Now().Add(-45 * time.Minute).UTC()
	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET status = 'FINISHED', finished_at = NOW(), last_activity_at = $2
		WHERE room_id = $1
	`, finishedRoomID, oldActivityAt); err != nil {
		t.Fatalf("failed to mark finished room: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), `
		UPDATE rooms
		SET last_activity_at = $2
		WHERE room_id = $1
	`, staleActiveRoomID, oldActivityAt); err != nil {
		t.Fatalf("failed to mark active room stale: %v", err)
	}

	expiryService := rooms.NewRoomsExpiryService(roomsrepositories.NewRoomsRepository(db), nil)
	expiredRooms, err := expiryService.ExpireInactiveRooms(time.Now().Add(-30 * time.Minute))
	if err != nil {
		t.Fatalf("failed to expire inactive rooms: %v", err)
	}
	if len(expiredRooms) != 1 {
		t.Fatalf("expected only the active stale room to expire, got %d", len(expiredRooms))
	}
	if expiredRooms[0].RoomID != staleActiveRoomID {
		t.Fatalf("expected stale active room %s to expire, got %s", staleActiveRoomID, expiredRooms[0].RoomID)
	}

	roomsRepo := roomsrepositories.NewRoomsRepository(db)
	finishedRoom, err := roomsRepo.FindByID(finishedRoomID)
	if err != nil {
		t.Fatalf("failed to load finished room: %v", err)
	}
	if finishedRoom.Status != "FINISHED" {
		t.Fatalf("expected finished room to stay FINISHED, got %s", finishedRoom.Status)
	}
	if finishedRoom.FinishedAt == nil {
		t.Fatal("expected finished room to keep finishedAt")
	}
}
