package roomsdto

import (
	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
)

type CreateRoomSkippedRecipientResponse struct {
	UserID *string `json:"userId,omitempty"`
	Email  *string `json:"email,omitempty"`
	Reason string  `json:"reason"`
}

type CreateRoomResponse struct {
	Room              *roomsmodels.RoomsModel                  `json:"room"`
	EmailInvites      []invitesdto.InvitationWithTokenResponse `json:"emailInvites,omitempty"`
	ShareLink         *invitesdto.InvitationWithTokenResponse  `json:"shareLink,omitempty"`
	InviteToken       string                                   `json:"inviteToken,omitempty"`
	SkippedRecipients []CreateRoomSkippedRecipientResponse     `json:"skippedRecipients,omitempty"`
}
