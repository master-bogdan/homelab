package roomsmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type RoomTaskModel struct {
	bun.BaseModel `bun:"table:tasks,alias:t"`

	TaskID             string    `bun:"task_id,pk"`
	RoomID             string    `bun:"room_id"`
	Title              string    `bun:"title"`
	Description        *string   `bun:"description"`
	ExternalKey        *string   `bun:"external_key"`
	Status             string    `bun:"status"`
	IsActive           bool      `bun:"is_active"`
	FinalEstimateValue *string   `bun:"final_estimate_value"`
	CreatedAt          time.Time `bun:"created_at"`
	UpdatedAt          time.Time `bun:"updated_at"`

	Votes []*RoomVoteModel `bun:"rel:has-many,join:task_id=task_id"`
	Room  *RoomsModel      `bun:"rel:belongs-to,join:room_id=room_id" json:"-"`
}
