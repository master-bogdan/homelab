package teamsmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type TeamModel struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	TeamID      string    `bun:"team_id,pk"`
	Name        string    `bun:"name"`
	OwnerUserID string    `bun:"owner_user_id"`
	CreatedAt   time.Time `bun:"created_at"`

	Members []*TeamMemberModel `bun:"rel:has-many,join:team_id=team_id"`
}
