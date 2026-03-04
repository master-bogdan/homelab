package roomsmodels

import (
	"time"

	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	"github.com/uptrace/bun"
)

type RoomParticipantRole string

const (
	RoomParticipantRoleAdmin  RoomParticipantRole = "ADMIN"
	RoomParticipantRoleMember RoomParticipantRole = "MEMBER"
	RoomParticipantRoleGuest  RoomParticipantRole = "GUEST"
)

func (r RoomParticipantRole) IsValid() bool {
	switch r {
	case RoomParticipantRoleAdmin, RoomParticipantRoleMember, RoomParticipantRoleGuest:
		return true
	default:
		return false
	}
}

type RoomParticipantModel struct {
	bun.BaseModel `bun:"table:room_participants,alias:rp"`

	RoomParticipantID string              `bun:"room_participants_id,pk"`
	RoomID            string              `bun:"room_id"`
	UserID            *string             `bun:"user_id"`
	GuestName         *string             `bun:"guest_name"`
	Role              RoomParticipantRole `bun:"role"`
	JoinedAt          time.Time           `bun:"joined_at"`
	LeftAt            *time.Time          `bun:"left_at"`

	User  *usersmodels.UserModel `bun:"rel:belongs-to,join:user_id=user_id"`
	Votes []*RoomVoteModel       `bun:"rel:has-many,join:room_participants_id=participant_id"`
	Room  *RoomsModel            `bun:"rel:belongs-to,join:room_id=room_id" json:"-"`
}
