package usersmodels

import "time"

type UserModel struct {
	UserID       string
	Email        *string
	PasswordHash *string
	GithubID     *string
	DisplayName  string
	AvatarURL    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
	DeletedAt    *time.Time
}
