// Package authrepositories includes auth repository implementations.
package authrepositories

import "errors"

var (
	ErrOidcSessionNotFound = errors.New("oidc session not found")
	ErrAccessTokenNotFound = errors.New("access token not found")
)
