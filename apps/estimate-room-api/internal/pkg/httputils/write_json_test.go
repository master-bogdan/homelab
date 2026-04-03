package httputils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

func TestWriteResponseErrorUsesRequestIDAsInstance(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.Header().Set(logger.RequestIDHeader, "req-123")

	WriteResponseError(rec, apperrors.CreateHttpError(
		apperrors.ErrUnauthorized,
		apperrors.HttpError{
			Detail:   "invalid credentials",
			Instance: "/api/v1/auth/login",
		},
	))

	var got apperrors.HttpError
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got.Instance != "req-123" {
		t.Fatalf("expected instance to be request id, got %q", got.Instance)
	}
}

func TestWriteResponseErrorKeepsExistingInstanceWithoutRequestID(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteResponseError(rec, apperrors.CreateHttpError(
		apperrors.ErrUnauthorized,
		apperrors.HttpError{
			Detail:   "invalid credentials",
			Instance: "/api/v1/auth/login",
		},
	))

	var got apperrors.HttpError
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got.Instance != "/api/v1/auth/login" {
		t.Fatalf("expected instance fallback to be preserved, got %q", got.Instance)
	}
}
