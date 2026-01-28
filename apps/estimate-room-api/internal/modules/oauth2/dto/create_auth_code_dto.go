package oauth2_dto

import (
	"errors"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CreateOauthCodeDTO struct {
	ClientID            string `validate:"required"`
	UserID              string `validate:"required"`
	OidcSessionID       string `validate:"required"`
	RedirectURI         string `validate:"required,url"`
	CodeChallenge       string `validate:"required"`
	CodeChallengeMethod string `validate:"required,oneof=S256"`
	Scopes              string `validate:"required"`
}

func (s *CreateOauthCodeDTO) Validate() error {
	validScopes := []string{"openid", "admin", "user"}

	scopes := strings.FieldsSeq(s.Scopes)

	for v := range scopes {
		isValidScope := slices.Contains(validScopes, v)

		if !isValidScope {
			return errors.New("no valid scope provided")
		}
	}

	validate := validator.New()

	return validate.Struct(s)
}
