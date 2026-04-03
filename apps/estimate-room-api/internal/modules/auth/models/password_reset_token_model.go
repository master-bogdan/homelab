package authmodels

import (
	"time"

	"github.com/uptrace/bun"
)

type PasswordResetTokenModel struct {
	bun.BaseModel `bun:"table:auth_password_reset_tokens,alias:prt"`

	PasswordResetTokenID string     `bun:"password_reset_token_id,pk"`
	UserID               string     `bun:"user_id"`
	TokenHash            string     `bun:"token_hash"`
	ExpiresAt            time.Time  `bun:"expires_at"`
	UsedAt               *time.Time `bun:"used_at"`
	CreatedAt            time.Time  `bun:"created_at"`
}
