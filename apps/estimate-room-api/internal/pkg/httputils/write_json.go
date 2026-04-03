// Package httputils contains HTTP response writers.
package httputils

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

type WriteResponseOptions struct {
	Status      int
	ContentType string
	Headers     map[string]string
}

func WriteResponse(w http.ResponseWriter, v any, opts ...WriteResponseOptions) {
	cfg := WriteResponseOptions{
		Status:      http.StatusOK,
		ContentType: "application/json; charset=utf-8",
	}
	if len(opts) > 0 {
		if opts[0].Status != 0 {
			cfg.Status = opts[0].Status
		}
		if opts[0].ContentType != "" {
			cfg.ContentType = opts[0].ContentType
		}
		if opts[0].Headers != nil {
			cfg.Headers = opts[0].Headers
		}
	}

	w.Header().Set("Content-Type", cfg.ContentType)
	for key, val := range cfg.Headers {
		w.Header().Set(key, val)
	}

	w.WriteHeader(cfg.Status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteResponseError(w http.ResponseWriter, httpErr apperrors.HttpError, opts ...WriteResponseOptions) {
	if httpErr.Errors == nil {
		httpErr.Errors = []apperrors.ErrorItem{}
	}

	cfg := WriteResponseOptions{
		Status:      httpErr.Status,
		ContentType: "application/problem+json; charset=utf-8",
	}
	if len(opts) > 0 {
		if opts[0].Status != 0 {
			cfg.Status = opts[0].Status
		}
		if opts[0].ContentType != "" {
			cfg.ContentType = opts[0].ContentType
		}
		if opts[0].Headers != nil {
			cfg.Headers = opts[0].Headers
		}
	}

	WriteResponse(w, httpErr, cfg)
}
