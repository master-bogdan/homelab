package roomsdto

import (
	"time"

	invitesdto "github.com/master-bogdan/estimate-room-api/internal/modules/invites/dto"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
)

type RoomDeckResponse struct {
	Name   string   `json:"name"`
	Kind   string   `json:"kind"`
	Values []string `json:"values"`
}

type RoomUserResponse struct {
	UserID      string  `json:"userId"`
	Email       *string `json:"email,omitempty"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type RoomParticipantResponse struct {
	RoomParticipantID string                          `json:"roomParticipantId"`
	RoomID            string                          `json:"roomId"`
	UserID            *string                         `json:"userId,omitempty"`
	GuestName         *string                         `json:"guestName,omitempty"`
	Role              roomsmodels.RoomParticipantRole `json:"role"`
	JoinedAt          time.Time                       `json:"joinedAt"`
	LeftAt            *time.Time                      `json:"leftAt,omitempty"`
	User              *RoomUserResponse               `json:"user,omitempty"`
}

type RoomTaskResponse struct {
	TaskID             string    `json:"taskId"`
	RoomID             string    `json:"roomId"`
	Title              string    `json:"title"`
	Description        *string   `json:"description,omitempty"`
	ExternalKey        *string   `json:"externalKey,omitempty"`
	Status             string    `json:"status"`
	IsActive           bool      `json:"isActive"`
	FinalEstimateValue *string   `json:"finalEstimateValue,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type RoomResponse struct {
	RoomID         string                    `json:"roomId"`
	Code           string                    `json:"code"`
	Name           string                    `json:"name"`
	AdminUserID    string                    `json:"adminUserId"`
	TeamID         *string                   `json:"teamId,omitempty"`
	Deck           RoomDeckResponse          `json:"deck"`
	Status         string                    `json:"status"`
	CreatedAt      time.Time                 `json:"createdAt"`
	LastActivityAt time.Time                 `json:"lastActivityAt"`
	FinishedAt     *time.Time                `json:"finishedAt,omitempty"`
	Participants   []RoomParticipantResponse `json:"participants,omitempty"`
	Tasks          []RoomTaskResponse        `json:"tasks,omitempty"`
}

type CreateRoomSkippedRecipientResponse struct {
	UserID *string `json:"userId,omitempty"`
	Email  *string `json:"email,omitempty"`
	Reason string  `json:"reason"`
}

type CreateRoomResponse struct {
	Room              *RoomResponse                            `json:"room"`
	EmailInvites      []invitesdto.InvitationWithTokenResponse `json:"emailInvites,omitempty"`
	ShareLink         *invitesdto.InvitationWithTokenResponse  `json:"shareLink,omitempty"`
	InviteToken       string                                   `json:"inviteToken,omitempty"`
	SkippedRecipients []CreateRoomSkippedRecipientResponse     `json:"skippedRecipients,omitempty"`
}

func NewRoomResponse(room *roomsmodels.RoomsModel) *RoomResponse {
	if room == nil {
		return nil
	}

	response := &RoomResponse{
		RoomID:         room.RoomID,
		Code:           room.Code,
		Name:           room.Name,
		AdminUserID:    room.AdminUserID,
		TeamID:         room.TeamID,
		Deck:           NewRoomDeckResponse(room.Deck),
		Status:         room.Status,
		CreatedAt:      room.CreatedAt,
		LastActivityAt: room.LastActivityAt,
		FinishedAt:     room.FinishedAt,
	}

	if len(room.Participants) > 0 {
		response.Participants = make([]RoomParticipantResponse, 0, len(room.Participants))
		for _, participant := range room.Participants {
			if participant == nil {
				continue
			}

			response.Participants = append(response.Participants, NewRoomParticipantResponse(participant))
		}
	}

	if len(room.Tasks) > 0 {
		response.Tasks = make([]RoomTaskResponse, 0, len(room.Tasks))
		for _, task := range room.Tasks {
			if task == nil {
				continue
			}

			response.Tasks = append(response.Tasks, NewRoomTaskResponse(task))
		}
	}

	return response
}

func NewRoomDeckResponse(deck roomsmodels.RoomDeck) RoomDeckResponse {
	return RoomDeckResponse{
		Name:   deck.Name,
		Kind:   deck.Kind,
		Values: deck.Values,
	}
}

func NewRoomParticipantResponse(participant *roomsmodels.RoomParticipantModel) RoomParticipantResponse {
	return RoomParticipantResponse{
		RoomParticipantID: participant.RoomParticipantID,
		RoomID:            participant.RoomID,
		UserID:            participant.UserID,
		GuestName:         participant.GuestName,
		Role:              participant.Role,
		JoinedAt:          participant.JoinedAt,
		LeftAt:            participant.LeftAt,
		User:              NewRoomUserResponse(participant.User),
	}
}

func NewRoomTaskResponse(task *roomsmodels.RoomTaskModel) RoomTaskResponse {
	return RoomTaskResponse{
		TaskID:             task.TaskID,
		RoomID:             task.RoomID,
		Title:              task.Title,
		Description:        task.Description,
		ExternalKey:        task.ExternalKey,
		Status:             task.Status,
		IsActive:           task.IsActive,
		FinalEstimateValue: task.FinalEstimateValue,
		CreatedAt:          task.CreatedAt,
		UpdatedAt:          task.UpdatedAt,
	}
}

func NewRoomUserResponse(user *usersmodels.UserModel) *RoomUserResponse {
	if user == nil {
		return nil
	}

	return &RoomUserResponse{
		UserID:      user.UserID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
	}
}
