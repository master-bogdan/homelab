package oauth2dto

import (
	"github.com/go-playground/validator/v10"
)

type UserDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

func (s *UserDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(s)
}
