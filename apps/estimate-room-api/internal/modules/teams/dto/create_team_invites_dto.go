package teamsdto

import "github.com/go-playground/validator/v10"

type CreateTeamInvitesDTO struct {
	Emails []string `json:"emails" validate:"required,min=1,max=200,dive,required,email,max=255"`
}

func (s *CreateTeamInvitesDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
