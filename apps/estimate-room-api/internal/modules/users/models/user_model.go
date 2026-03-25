package usersmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	UserID       string     `bun:"user_id,pk"`
	Email        *string    `bun:"email"`
	PasswordHash *string    `bun:"password_hash"`
	GithubID     *string    `bun:"github_id"`
	DisplayName  string     `bun:"display_name"`
	Organization *string    `bun:"organization"`
	Occupation   *string    `bun:"occupation"`
	AvatarURL    *string    `bun:"avatar_url"`
	CreatedAt    time.Time  `bun:"created_at"`
	UpdatedAt    time.Time  `bun:"updated_at"`
	LastLoginAt  *time.Time `bun:"last_login_at"`
	DeletedAt    *time.Time `bun:"deleted_at"`
}
