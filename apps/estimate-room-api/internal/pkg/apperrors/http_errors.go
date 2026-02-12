package apperrors

import (
	stdErrors "errors"
	"net/http"
)

var (
	ErrBadRequest         = stdErrors.New("bad request")
	ErrUnauthorized       = stdErrors.New("unauthorized")
	ErrForbidden          = stdErrors.New("forbidden")
	ErrNotFound           = stdErrors.New("not found")
	ErrConflict           = stdErrors.New("conflict")
	ErrInternal           = stdErrors.New("internal server error")
	ErrServiceUnavailable = stdErrors.New("service unavailable")
)

type HttpError struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   string      `json:"detail"`
	Instance string      `json:"instance"`
	Errors   []ErrorItem `json:"errors"`
}

func CreateHttpError(errType error, overrides ...HttpError) HttpError {
	httpErr := defaultHttpError(errType)

	if len(overrides) > 0 {
		override := overrides[0]
		if override.Type != "" {
			httpErr.Type = override.Type
		}
		if override.Title != "" {
			httpErr.Title = override.Title
		}
		if override.Status != 0 {
			httpErr.Status = override.Status
		}
		if override.Detail != "" {
			httpErr.Detail = override.Detail
		}
		if override.Instance != "" {
			httpErr.Instance = override.Instance
		}
		if override.Errors != nil {
			httpErr.Errors = override.Errors
		}
	}

	if httpErr.Errors == nil {
		httpErr.Errors = []ErrorItem{}
	}

	return httpErr
}

func defaultHttpError(errType error) HttpError {
	switch {
	case stdErrors.Is(errType, ErrBadRequest):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/bad-request",
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "bad request",
		}
	case stdErrors.Is(errType, ErrUnauthorized):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/unauthorized",
			Title:  "Unauthorized",
			Status: http.StatusUnauthorized,
			Detail: "unauthorized",
		}
	case stdErrors.Is(errType, ErrForbidden):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/forbidden",
			Title:  "Forbidden",
			Status: http.StatusForbidden,
			Detail: "forbidden",
		}
	case stdErrors.Is(errType, ErrUserNotFound):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/not-found",
			Title:  "Not Found",
			Status: http.StatusNotFound,
			Detail: "user not found",
		}
	case stdErrors.Is(errType, ErrNotFound):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/not-found",
			Title:  "Not Found",
			Status: http.StatusNotFound,
			Detail: "resource not found",
		}
	case stdErrors.Is(errType, ErrConflict):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/conflict",
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: "conflict",
		}
	case stdErrors.Is(errType, ErrServiceUnavailable):
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/service-unavailable",
			Title:  "Service Unavailable",
			Status: http.StatusServiceUnavailable,
			Detail: "service unavailable",
		}
	default:
		return HttpError{
			Type:   "https://api.estimateroom.com/problems/internal-server-error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "internal server error",
		}
	}
}
