package authdto

import "github.com/go-playground/validator/v10"

type RegisterDTO struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Password    string `json:"password" validate:"required,min=8,max=128"`
	DisplayName string `json:"displayName" validate:"omitempty,min=1,max=100"`
	ContinueURL string `json:"continue" validate:"required,url,max=2000"`
}

func (s *RegisterDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
