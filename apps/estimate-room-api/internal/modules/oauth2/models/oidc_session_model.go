package oauth2models

import (
	"time"

	"github.com/uptrace/bun"
)

type OidcSessionModel struct {
	bun.BaseModel `bun:"table:oauth2_oidc_sessions,alias:os"`

	OidcSessionID string     `bun:"oidc_session_id,pk"`
	UserID        string     `bun:"user_id"`
	ClientID      string     `bun:"client_id"`
	Nonce         string     `bun:"nonce"`
	CreatedAt     time.Time  `bun:"created_at"`
	RevokedAt     *time.Time `bun:"revoked_at"`

	Client        *Oauth2ClientModel         `bun:"rel:belongs-to,join:client_id=client_id"`
	AuthCodes     []*Oauth2AuthCodeModel     `bun:"rel:has-many,join:oidc_session_id=oidc_session_id"`
	RefreshTokens []*Oauth2RefreshTokenModel `bun:"rel:has-many,join:oidc_session_id=oidc_session_id"`
	AccessTokens  []*Oauth2AccessTokenModel  `bun:"rel:has-many,join:oidc_session_id=oidc_session_id"`
}
