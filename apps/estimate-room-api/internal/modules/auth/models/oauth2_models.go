package authmodels

import "time"

type Oauth2AccessTokenModel struct {
	AccessTokenID  string
	UserID         string
	ClientID       string
	OidcSessionID  string
	RefreshTokenID *string
	Scopes         []string
	Token          string
	IssuedAt       time.Time
	ExpiresAt      time.Time
	Issuer         string
	IsRevoked      bool
	CreatedAt      time.Time
}

type OidcSessionModel struct {
	OidcSessionID string
	UserID        string
	ClientID      string
	Nonce         string
	CreatedAt     time.Time
}
