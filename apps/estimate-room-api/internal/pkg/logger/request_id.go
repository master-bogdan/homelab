package logger

import (
	"net/http"

	"github.com/google/uuid"
)

// RequestIDHeader is the HTTP header used to propagate request IDs.
const RequestIDHeader = "X-Request-Id"

// RequestIDMiddleware ensures every request has a request ID in context and response.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		ctx := WithRequestID(r.Context(), requestID)
		w.Header().Set(RequestIDHeader, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
