package teamsdto

import "github.com/go-playground/validator/v10"

type CreateTeamDTO struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (s *CreateTeamDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
