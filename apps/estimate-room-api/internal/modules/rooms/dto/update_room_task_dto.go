package roomsdto

import (
	"github.com/go-playground/validator/v10"
)

type UpdateRoomTaskDTO struct {
	Title              *string `json:"title" validate:"omitempty,min=1,max=255"`
	Description        *string `json:"description" validate:"omitempty,max=2000"`
	ExternalKey        *string `json:"externalKey" validate:"omitempty,max=255"`
	Status             *string `json:"status" validate:"omitempty,oneof=PENDING VOTING ESTIMATED SKIPPED"`
	IsActive           *bool   `json:"isActive"`
	FinalEstimateValue *string `json:"finalEstimateValue" validate:"omitempty,max=255"`
}

func (s *UpdateRoomTaskDTO) Validate() error {
	validate := validator.New()

	return validate.Struct(s)
}
