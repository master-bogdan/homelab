package invitesmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type InvitationKind string

const (
	InvitationKindTeamMember InvitationKind = "TEAM_MEMBER"
	InvitationKindRoomEmail  InvitationKind = "ROOM_EMAIL"
	InvitationKindRoomLink   InvitationKind = "ROOM_LINK"
)

func (k InvitationKind) IsValid() bool {
	switch k {
	case InvitationKindTeamMember, InvitationKindRoomEmail, InvitationKindRoomLink:
		return true
	default:
		return false
	}
}

type InvitationStatus string

const (
	InvitationStatusActive   InvitationStatus = "ACTIVE"
	InvitationStatusAccepted InvitationStatus = "ACCEPTED"
	InvitationStatusDeclined InvitationStatus = "DECLINED"
	InvitationStatusRevoked  InvitationStatus = "REVOKED"
)

func (s InvitationStatus) IsValid() bool {
	switch s {
	case InvitationStatusActive, InvitationStatusAccepted, InvitationStatusDeclined, InvitationStatusRevoked:
		return true
	default:
		return false
	}
}

type InvitationModel struct {
	bun.BaseModel `bun:"table:invitations,alias:i"`

	InvitationID    string           `bun:"invitation_id,pk"`
	Kind            InvitationKind   `bun:"kind"`
	Status          InvitationStatus `bun:"status"`
	TeamID          *string          `bun:"team_id"`
	RoomID          *string          `bun:"room_id"`
	InvitedUserID   *string          `bun:"invited_user_id"`
	InvitedEmail    *string          `bun:"invited_email"`
	CreatedByUserID string           `bun:"created_by_user_id"`
	TokenID         string           `bun:"token_id"`
	AcceptedAt      *time.Time       `bun:"accepted_at"`
	DeclinedAt      *time.Time       `bun:"declined_at"`
	RevokedAt       *time.Time       `bun:"revoked_at"`
	CreatedAt       time.Time        `bun:"created_at"`
	UpdatedAt       time.Time        `bun:"updated_at"`
}
