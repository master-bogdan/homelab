// Package roomsdto is a collection of rooms dtos
package roomsdto

import "github.com/go-playground/validator/v10"

type CreateRoomDTO struct {
	Name   string             `json:"name" validate:"required,min=1,max=30"`
	TeamID string             `json:"teamId" validate:"omitempty"`
	Deck   *CreateRoomDeckDTO `json:"deck" validate:"omitempty"`
}

type CreateRoomDeckDTO struct {
	Name   string   `json:"name" validate:"required,min=1,max=50"`
	Kind   string   `json:"kind" validate:"required,min=1,max=30"`
	Values []string `json:"values" validate:"required,min=1,max=50,dive,required,max=20"`
}

func (s *CreateRoomDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
