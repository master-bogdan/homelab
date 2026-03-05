package roomsmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type RoomTaskRoundModel struct {
	bun.BaseModel `bun:"table:task_rounds,alias:tr"`

	TaskID      string     `bun:"task_id,pk"`
	RoundNumber int        `bun:"round_number,pk"`
	IsRevealed  bool       `bun:"is_revealed"`
	RevealedAt  *time.Time `bun:"revealed_at"`
	CreatedAt   time.Time  `bun:"created_at"`
	UpdatedAt   time.Time  `bun:"updated_at"`

	Task *RoomTaskModel `bun:"rel:belongs-to,join:task_id=task_id" json:"-"`
}
