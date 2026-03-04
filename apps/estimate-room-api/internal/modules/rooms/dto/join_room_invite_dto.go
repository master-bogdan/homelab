package roomsdto

import "github.com/go-playground/validator/v10"

type JoinRoomInviteDTO struct {
	GuestName *string `json:"guestName" validate:"omitempty,min=1,max=100"`
}

func (s *JoinRoomInviteDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
