package oauth2_dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type GetTokenDTO struct {
	GrantType    string `form:"grant_type" validate:"required,oneof=authorization_code refresh_token"`
	CodeVerifier string `form:"code_verifier"`
	Code         string `form:"code"`
	ClientID     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	RefreshToken string `form:"refresh_token"`
}

func (s *GetTokenDTO) Validate() error {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return err
	}

	switch s.GrantType {
	case "authorization_code":
		if s.Code == "" || s.CodeVerifier == "" || s.ClientID == "" || s.RedirectURI == "" {
			return errors.New("code, code_verifier, client_id, and redirect_uri are required for authorization_code")
		}
	case "refresh_token":
		if s.RefreshToken == "" || s.ClientID == "" {
			return errors.New("refresh_token and client_id are required for refresh_token")
		}
	default:
		return errors.New("unsupported grant_type")
	}

	return nil
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
