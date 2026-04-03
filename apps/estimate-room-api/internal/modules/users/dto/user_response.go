package usersdto

import "time"

type UserResponse struct {
	ID          string     `json:"id"`
	Email       *string    `json:"email,omitempty"`
	GithubID    *string    `json:"githubId,omitempty"`
	DisplayName string     `json:"displayName"`
	AvatarURL   *string    `json:"avatarUrl,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
