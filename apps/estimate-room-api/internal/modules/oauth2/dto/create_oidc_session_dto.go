package oauth2_dto

import "github.com/go-playground/validator/v10"

type CreateOidcSessionDTO struct {
	ClientID string `validate:"required"`
	UserID   string `validate:"required"`
	Nonce    string `validate:"required"`
}

func (s *CreateOidcSessionDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(s)
}
