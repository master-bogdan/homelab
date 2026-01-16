package oauth2_dto

import (
	"errors"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AuthorizeQueryDTO struct {
	ClientID            string `query:"client_id" validate:"required"`
	RedirectURI         string `query:"redirect_uri" validate:"required,url"`
	ResponseType        string `query:"response_type" validate:"required,eq=code"`
	Scopes              string `query:"scopes" validate:"required"`
	State               string `query:"state" validate:"required"`
	CodeChallenge       string `query:"code_challenge" validate:"required"`
	CodeChallengeMethod string `query:"code_challenge_method" validate:"required"`
	Nonce               string `query:"nonce" validate:"required"`
}

func (s *AuthorizeQueryDTO) Validate() error {
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
