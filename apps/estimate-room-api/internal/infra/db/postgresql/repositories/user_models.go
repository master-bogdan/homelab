package repositories

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

type UserRepository interface {
	FindByID(userID string) (*UserModel, error)
	FindByEmail(email string) (*UserModel, error)
	Create(email, passwordHash string) (string, error)
}
