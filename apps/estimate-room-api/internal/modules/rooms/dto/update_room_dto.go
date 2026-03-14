package roomsdto

import (
	"github.com/go-playground/validator/v10"
)

type UpdateRoomDTO struct {
	Name   *string `json:"name" validate:"omitempty,min=1,max=30"`
	Status *string `json:"status" validate:"omitempty,oneof=ACTIVE FINISHED EXPIRED"`
}

func (s *UpdateRoomDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
