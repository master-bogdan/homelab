// Package oauth2repositories includes oauth2 repository implementations.
package oauth2repositories

import "errors"

var (
	ErrClientNotFound       = errors.New("client not found")
	ErrAuthCodeNotFound     = errors.New("auth code not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)
