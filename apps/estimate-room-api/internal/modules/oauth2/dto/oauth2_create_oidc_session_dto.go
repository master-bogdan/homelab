package oauth2dto

import "github.com/go-playground/validator/v10"

type Oauth2CreateOidcSessionDTO struct {
	ClientID string `validate:"required"`
	UserID   string `validate:"required"`
	Nonce    string `validate:"required"`
}

func (s *Oauth2CreateOidcSessionDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(s)
}
