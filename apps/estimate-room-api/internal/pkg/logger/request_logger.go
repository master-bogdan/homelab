package logger

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/metrics"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

func (r *responseRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not support hijacking")
	}
	return h.Hijack()
}

func (r *responseRecorder) Push(target string, opts *http.PushOptions) error {
	if p, ok := r.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w}

		L().InfoContext(r.Context(), Prefix("MW", "HTTP", "REQUEST", "Request started"),
			"method", r.Method,
			"path", r.URL.Path,
			"headers", requestHeaders(r),
			"content_length", r.ContentLength,
			"remote_ip", clientIP(r),
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(rec, r)

		status := rec.status
		if status == 0 {
			status = http.StatusOK
		}

		level := slog.LevelInfo
		if status >= http.StatusInternalServerError {
			level = slog.LevelError
		} else if status >= http.StatusBadRequest {
			level = slog.LevelWarn
		}

		duration := time.Since(start)
		routePattern := r.URL.Path
		if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil {
			if pattern := strings.TrimSpace(routeCtx.RoutePattern()); pattern != "" {
				routePattern = pattern
			}
		}
		metrics.RecordHTTPRequest(r.Method, routePattern, status, duration)

		L().Log(r.Context(), level, Prefix("MW", "HTTP", "REQUEST", "Request ended"),
			"method", r.Method,
			"path", r.URL.Path,
			"route", routePattern,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"bytes", rec.bytes,
			"remote_ip", clientIP(r),
			"user_agent", r.UserAgent(),
		)
	})
}

func requestHeaders(r *http.Request) map[string]string {
	if r == nil {
		return map[string]string{}
	}

	allowedHeaders := []string{
		"Accept",
		"Content-Type",
		"Origin",
		"Referer",
		RequestIDHeader,
	}

	headers := make(map[string]string, len(allowedHeaders))
	for _, header := range allowedHeaders {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}

		headers[strings.ToLower(header)] = value
	}

	if len(headers) == 0 {
		return map[string]string{}
	}

	normalizedHeaders := make(map[string]string, len(headers))
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		normalizedHeaders[key] = headers[key]
	}

	return normalizedHeaders
}

func clientIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}
