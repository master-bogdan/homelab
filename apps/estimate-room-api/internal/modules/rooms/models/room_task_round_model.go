package roomsmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type RoomTaskRoundStatus string

const (
	RoomTaskRoundStatusActive   RoomTaskRoundStatus = "ACTIVE"
	RoomTaskRoundStatusRevealed RoomTaskRoundStatus = "REVEALED"
)

type RoomTaskRoundModel struct {
	bun.BaseModel `bun:"table:task_rounds,alias:tr"`

	TaskID                 string              `bun:"task_id,pk"`
	RoundNumber            int                 `bun:"round_number,pk"`
	EligibleParticipantIDs []string            `bun:"eligible_participant_ids,type:jsonb"`
	Status                 RoomTaskRoundStatus `bun:"status"`
	CreatedAt              time.Time           `bun:"created_at"`
	UpdatedAt              time.Time           `bun:"updated_at"`

	Task *RoomTaskModel `bun:"rel:belongs-to,join:task_id=task_id" json:"-"`
}
