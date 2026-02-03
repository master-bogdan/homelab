package models

import "time"

type Oauth2ClientModel struct {
	ClientID      string
	ClientSecret  string
	RedirectURIs  []string
	GrantTypes    []string
	ResponseTypes []string
	Scopes        []string
	ClientName    string
	ClientType    string
	CreatedAt     time.Time
}

type Oauth2AuthCodeModel struct {
	AuthCodeID          string
	ClientID            string
	UserID              string
	OidcSessionID       string
	Code                string
	RedirectURI         string
	Scopes              []string
	CodeChallenge       string
	CodeChallengeMethod string
	IsUsed              bool
	ExpiresAt           time.Time
	CreatedAt           time.Time
}

type Oauth2RefreshTokenModel struct {
	RefreshTokenID string
	UserID         string
	ClientID       string
	OidcSessionID  string
	Scopes         []string
	Token          string
	IssuedAt       time.Time
	ExpiresAt      time.Time
	IsRevoked      bool
	CreatedAt      time.Time
}

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
