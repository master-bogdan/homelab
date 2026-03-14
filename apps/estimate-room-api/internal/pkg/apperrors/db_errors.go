package apperrors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrClientNotFound       = errors.New("client not found")
	ErrAuthCodeNotFound     = errors.New("auth code not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrOidcSessionNotFound  = errors.New("oidc session not found")
	ErrAccessTokenNotFound  = errors.New("access token not found")
)
