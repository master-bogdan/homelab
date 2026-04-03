package teamsmodels

import (
	"time"

	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	"github.com/uptrace/bun"
)

type TeamMemberRole string

const (
	TeamMemberRoleOwner  TeamMemberRole = "OWNER"
	TeamMemberRoleMember TeamMemberRole = "MEMBER"
)

type TeamMemberModel struct {
	bun.BaseModel `bun:"table:team_members,alias:tm"`

	TeamID   string         `bun:"team_id,pk"`
	UserID   string         `bun:"user_id,pk"`
	Role     TeamMemberRole `bun:"role"`
	JoinedAt time.Time      `bun:"joined_at"`

	User *usersmodels.UserModel `bun:"rel:belongs-to,join:user_id=user_id"`
	Team *TeamModel             `bun:"rel:belongs-to,join:team_id=team_id" json:"-"`
}
