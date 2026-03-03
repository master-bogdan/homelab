package roomsdto

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type UpdateRoomDTO struct {
	Name              *string `json:"name" validate:"omitempty,min=1,max=30"`
	Status            *string `json:"status" validate:"omitempty,oneof=ACTIVE FINISHED EXPIRED"`
	AllowGuests       *bool   `json:"allowGuests"`
	AllowSpectators   *bool   `json:"allowSpectators"`
	RoundTimerSeconds *int    `json:"roundTimerSeconds" validate:"omitempty,min=1,max=86400"`
}

func (s *UpdateRoomDTO) Validate() error {
	if s.Name == nil &&
		s.Status == nil &&
		s.AllowGuests == nil &&
		s.AllowSpectators == nil &&
		s.RoundTimerSeconds == nil {
		return errors.New("at least one field must be provided")
	}

	validate := validator.New()
	return validate.Struct(s)
}
