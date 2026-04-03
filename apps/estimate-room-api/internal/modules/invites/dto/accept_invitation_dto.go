package invitesdto

import "github.com/go-playground/validator/v10"

type AcceptInvitationDTO struct {
	GuestName *string `json:"guestName" validate:"omitempty,min=1,max=100"`
}

func (s *AcceptInvitationDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
