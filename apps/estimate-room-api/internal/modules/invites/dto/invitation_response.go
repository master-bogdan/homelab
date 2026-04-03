package invitesdto

import (
	"time"

	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
)

type InvitationResponse struct {
	InvitationID    string     `json:"invitationId"`
	Kind            string     `json:"kind"`
	Status          string     `json:"status"`
	TeamID          *string    `json:"teamId"`
	RoomID          *string    `json:"roomId"`
	InvitedUserID   *string    `json:"invitedUserId"`
	InvitedEmail    *string    `json:"invitedEmail"`
	CreatedByUserID string     `json:"createdByUserId"`
	AcceptedAt      *time.Time `json:"acceptedAt"`
	DeclinedAt      *time.Time `json:"declinedAt"`
	RevokedAt       *time.Time `json:"revokedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type InvitationWithTokenResponse struct {
	InvitationResponse
	Token string `json:"token"`
}

func NewInvitationResponse(invitation *invitesmodels.InvitationModel) InvitationResponse {
	return InvitationResponse{
		InvitationID:    invitation.InvitationID,
		Kind:            string(invitation.Kind),
		Status:          string(invitation.Status),
		TeamID:          invitation.TeamID,
		RoomID:          invitation.RoomID,
		InvitedUserID:   invitation.InvitedUserID,
		InvitedEmail:    invitation.InvitedEmail,
		CreatedByUserID: invitation.CreatedByUserID,
		AcceptedAt:      invitation.AcceptedAt,
		DeclinedAt:      invitation.DeclinedAt,
		RevokedAt:       invitation.RevokedAt,
		CreatedAt:       invitation.CreatedAt,
		UpdatedAt:       invitation.UpdatedAt,
	}
}

func NewInvitationWithTokenResponse(
	invitation *invitesmodels.InvitationModel,
	token string,
) InvitationWithTokenResponse {
	return InvitationWithTokenResponse{
		InvitationResponse: NewInvitationResponse(invitation),
		Token:              token,
	}
}
