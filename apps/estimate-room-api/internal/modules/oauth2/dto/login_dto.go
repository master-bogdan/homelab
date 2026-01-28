package oauth2dto

import (
	"errors"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type LoginDTO struct {
	Email               string `form:"email" validate:"required,email"`
	Password            string `form:"password" validate:"required,min=6"`
	ClientID            string `form:"client_id" validate:"required"`
	RedirectURI         string `form:"redirect_uri" validate:"required,url"`
	ResponseType        string `form:"response_type" validate:"required,eq=code"`
	Scopes              string `form:"scopes" validate:"required"`
	State               string `form:"state" validate:"required"`
	CodeChallenge       string `form:"code_challenge" validate:"required"`
	CodeChallengeMethod string `form:"code_challenge_method" validate:"required,oneof=S256"`
	Nonce               string `form:"nonce" validate:"required"`
}

func (s *LoginDTO) Validate() error {
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
