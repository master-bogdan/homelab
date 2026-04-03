package oauth2models

import (
	"time"

	"github.com/uptrace/bun"
)

type Oauth2AuthCodeModel struct {
	bun.BaseModel `bun:"table:oauth2_auth_codes,alias:oac"`

	AuthCodeID          string    `bun:"auth_code_id,pk"`
	ClientID            string    `bun:"client_id"`
	UserID              string    `bun:"user_id"`
	OidcSessionID       string    `bun:"oidc_session_id"`
	Code                string    `bun:"code"`
	RedirectURI         string    `bun:"redirect_uri"`
	Scopes              []string  `bun:"scopes"`
	CodeChallenge       string    `bun:"code_challenge"`
	CodeChallengeMethod string    `bun:"code_challenge_method"`
	IsUsed              bool      `bun:"is_used"`
	ExpiresAt           time.Time `bun:"expires_at"`
	CreatedAt           time.Time `bun:"created_at"`

	Client      *Oauth2ClientModel `bun:"rel:belongs-to,join:client_id=client_id"`
	OidcSession *OidcSessionModel  `bun:"rel:belongs-to,join:oidc_session_id=oidc_session_id"`
}
