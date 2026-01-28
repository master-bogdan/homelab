package utils

import (
	"crypto/sha256"
	"errors"

	"github.com/o1egl/paseto"
)

func normalizeTokenKey(key []byte) []byte {
	if len(key) == 32 {
		return key
	}
	sum := sha256.Sum256(key)
	return sum[:]
}

func GenerateToken[T any](key []byte, payload T) (string, error) {
	if len(key) == 0 {
		return "", errors.New("token key is required")
	}

	normalized := normalizeTokenKey(key)
	return paseto.NewV2().Encrypt(normalized, payload, nil)
}

func ParseToken[T any](key []byte, token string) (*T, error) {
	if len(key) == 0 {
		return nil, errors.New("token key is required")
	}

	var payload T
	normalized := normalizeTokenKey(key)
	if err := paseto.NewV2().Decrypt(token, normalized, &payload, nil); err != nil {
		return nil, err
	}

	return &payload, nil
}
