package roomsmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type RoomsModel struct {
	bun.BaseModel `bun:"table:rooms,alias:r"`

	RoomID         string     `bun:"room_id,pk"`
	Code           string     `bun:"code"`
	Name           string     `bun:"name"`
	AdminUserID    string     `bun:"admin_user_id"`
	TeamID         *string    `bun:"team_id"`
	DeckID         DeckID     `bun:"deck_id"`
	Status         string     `bun:"status"`
	CreatedAt      time.Time  `bun:"created_at"`
	LastActivityAt time.Time  `bun:"last_activity_at"`
	FinishedAt     *time.Time `bun:"finished_at"`

	Participants []*RoomParticipantModel `bun:"rel:has-many,join:room_id=room_id"`
	Tasks        []*RoomTaskModel        `bun:"rel:has-many,join:room_id=room_id"`
}
