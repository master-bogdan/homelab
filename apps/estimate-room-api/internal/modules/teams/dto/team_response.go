package teamsdto

import (
	"time"

	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
)

type TeamUserResponse struct {
	UserID      string  `json:"userId"`
	Email       *string `json:"email"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl"`
}

type TeamMemberResponse struct {
	UserID   string           `json:"userId"`
	Role     string           `json:"role"`
	JoinedAt time.Time        `json:"joinedAt"`
	User     TeamUserResponse `json:"user"`
}

type TeamSummaryResponse struct {
	TeamID      string    `json:"teamId"`
	Name        string    `json:"name"`
	OwnerUserID string    `json:"ownerUserId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TeamDetailResponse struct {
	TeamSummaryResponse
	Members []TeamMemberResponse `json:"members"`
}

func NewTeamSummaryResponse(team *teamsmodels.TeamModel) TeamSummaryResponse {
	return TeamSummaryResponse{
		TeamID:      team.TeamID,
		Name:        team.Name,
		OwnerUserID: team.OwnerUserID,
		CreatedAt:   team.CreatedAt,
	}
}

func NewTeamDetailResponse(team *teamsmodels.TeamModel) TeamDetailResponse {
	members := make([]TeamMemberResponse, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, NewTeamMemberResponse(member))
	}

	return TeamDetailResponse{
		TeamSummaryResponse: NewTeamSummaryResponse(team),
		Members:             members,
	}
}

func NewTeamMemberResponse(member *teamsmodels.TeamMemberModel) TeamMemberResponse {
	response := TeamMemberResponse{
		UserID:   member.UserID,
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt,
	}

	if member.User != nil {
		response.User = TeamUserResponse{
			UserID:      member.User.UserID,
			Email:       member.User.Email,
			DisplayName: member.User.DisplayName,
			AvatarURL:   member.User.AvatarURL,
		}
	} else {
		response.User = TeamUserResponse{
			UserID: member.UserID,
		}
	}

	return response
}
