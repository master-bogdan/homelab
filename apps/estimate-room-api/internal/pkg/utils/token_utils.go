package utils

import (
	"github.com/o1egl/paseto"
)

func GenerateToken[T any](key []byte, payload T) (string, error) {
	return paseto.NewV2().Encrypt(key, payload, nil)
}

func ParseToken[T any](key []byte, token string) (*T, error) {
	var payload T
	err := paseto.NewV2().Decrypt(token, key, &payload, nil)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
