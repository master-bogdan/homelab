package oauth2_dto

import (
	"github.com/go-playground/validator/v10"
)

type GetTokenDTO struct {
	GrantType    string `form:"grant_type" validate:"required,oneof=authorization_code refresh_token"`
	CodeVerifier string `form:"code_verifier" validate:"required"`
	Code         string `form:"code"`
	ClientID     string `form:"client_id"`
	RefreshToken string `form:"refresh_token"`
}

// TODO: add better validation depending on grant_type

func (s *GetTokenDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(s)
}

type IDTokenPayload struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	AuthTime  int64  `json:"auth_time,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"` // "Bearer"
	ExpiresIn    int    `json:"expires_in"` // in seconds
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}
