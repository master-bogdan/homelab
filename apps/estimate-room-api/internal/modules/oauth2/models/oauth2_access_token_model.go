package oauth2models

import (
	"time"

	"github.com/uptrace/bun"
)

type Oauth2AccessTokenModel struct {
	bun.BaseModel `bun:"table:oauth2_access_tokens,alias:oat"`

	AccessTokenID  string     `bun:"access_token_id,pk"`
	UserID         string     `bun:"user_id"`
	ClientID       string     `bun:"client_id"`
	OidcSessionID  string     `bun:"oidc_session_id"`
	RefreshTokenID *string    `bun:"refresh_token_id"`
	Scopes         []string   `bun:"scopes"`
	Token          string     `bun:"token"`
	IssuedAt       time.Time  `bun:"issued_at"`
	ExpiresAt      time.Time  `bun:"expires_at"`
	Issuer         string     `bun:"issuer"`
	IsRevoked      bool       `bun:"is_revoked"`
	CreatedAt      time.Time  `bun:"created_at"`

	Client       *Oauth2ClientModel       `bun:"rel:belongs-to,join:client_id=client_id"`
	OidcSession  *OidcSessionModel        `bun:"rel:belongs-to,join:oidc_session_id=oidc_session_id"`
	RefreshToken *Oauth2RefreshTokenModel `bun:"rel:belongs-to,join:refresh_token_id=refresh_token_id"`
}
