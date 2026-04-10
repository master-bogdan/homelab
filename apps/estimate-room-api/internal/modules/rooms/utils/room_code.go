package roomsutils

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/oklog/ulid/v2"
)

func GenerateRoomCode() (string, error) {
	code, err := ulid.New(ulid.Timestamp(time.Now().UTC()), rand.Reader)
	if err != nil {
		return "", fmt.Errorf("%w: failed to generate room code", apperrors.ErrInternal)
	}

	return code.String(), nil
}
