package roomsmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type RoomVoteModel struct {
	bun.BaseModel `bun:"table:votes,alias:v"`

	VoteID        string    `bun:"votes_id,pk"`
	TaskID        string    `bun:"task_id"`
	ParticipantID string    `bun:"participant_id"`
	Value         string    `bun:"value"`
	RoundNumber   int       `bun:"round_number"`
	CreatedAt     time.Time `bun:"created_at"`

	Task        *RoomTaskModel        `bun:"rel:belongs-to,join:task_id=task_id" json:"-"`
	Participant *RoomParticipantModel `bun:"rel:belongs-to,join:participant_id=room_participants_id" json:"-"`
}
