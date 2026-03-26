package oauth2models

import (
	"time"

	"github.com/uptrace/bun"
)

type Oauth2ClientModel struct {
	bun.BaseModel `bun:"table:oauth2_clients,alias:oc"`

	ClientID      string    `bun:"client_id,pk"`
	ClientSecret  string    `bun:"client_secret"`
	RedirectURIs  []string  `bun:"redirect_uris,array"`
	GrantTypes    []string  `bun:"grant_types,array"`
	ResponseTypes []string  `bun:"response_types,array"`
	Scopes        []string  `bun:"scopes,array"`
	ClientName    string    `bun:"client_name"`
	ClientType    string    `bun:"client_type"`
	CreatedAt     time.Time `bun:"created_at"`

	AuthCodes     []*Oauth2AuthCodeModel     `bun:"rel:has-many,join:client_id=client_id"`
	RefreshTokens []*Oauth2RefreshTokenModel `bun:"rel:has-many,join:client_id=client_id"`
	AccessTokens  []*Oauth2AccessTokenModel  `bun:"rel:has-many,join:client_id=client_id"`
	OidcSessions  []*OidcSessionModel        `bun:"rel:has-many,join:client_id=client_id"`
}
