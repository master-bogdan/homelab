package roomsdto

import "github.com/go-playground/validator/v10"

type CreateRoomTaskDTO struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=2000"`
	ExternalKey string `json:"externalKey" validate:"omitempty,max=255"`
}

func (s *CreateRoomTaskDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
