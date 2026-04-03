package authdto

import "github.com/go-playground/validator/v10"

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (s *ForgotPasswordDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
