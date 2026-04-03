package authdto

import "github.com/go-playground/validator/v10"

type ResetPasswordDTO struct {
	Token    string `json:"token" validate:"required,max=512"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

func (s *ResetPasswordDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
