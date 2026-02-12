// Package roomsdto is a collection of rooms dtos
package roomsdto

import "github.com/go-playground/validator/v10"

type CreateRoomDTO struct {
	Name   string `json:"name" validate:"required, min=1, max=30"`
	TeamID string `json:"teamId" validate:"omitempty"`
	DeckID string `json:"deckId" validate:"omitempty,oneof=FIBONACCI TSHIRT CUSTOM"`
}

func (s *CreateRoomDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
