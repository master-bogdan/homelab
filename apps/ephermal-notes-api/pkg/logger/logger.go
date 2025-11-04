package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

type JSONHandler struct {
	level  slog.Leveler
	writer *os.File
	color  bool
	attrs  []slog.Attr
}

func NewJSONHandler(level slog.Leveler) *JSONHandler {
	return &JSONHandler{
		level:  level,
		writer: os.Stdout,
		color:  isTerminal(os.Stdout),
	}
}

func (h *JSONHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level.Level()
}

func (h *JSONHandler) Handle(ctx context.Context, r slog.Record) error {
	logEntry := map[string]any{
		"@timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"level":      r.Level.String(),
		"message":    r.Message,
	}

	if rid, ok := ctx.Value(requestIDKey).(string); ok && rid != "" {
		logEntry["request_id"] = rid
	}

	for _, a := range h.attrs {
		logEntry[a.Key] = a.Value.Any()
	}
	r.Attrs(func(a slog.Attr) bool {
		logEntry[a.Key] = a.Value.Any()
		return true
	})

	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	if h.color {
		colored := colorize(r.Level, string(data))
		_, err = fmt.Fprintln(h.writer, colored)
		return err
	}

	_, err = h.writer.Write(append(data, '\n'))
	return err
}

func (h *JSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &JSONHandler{
		level:  h.level,
		writer: h.writer,
		color:  h.color,
		attrs:  append(h.attrs, attrs...),
	}
}

func (h *JSONHandler) WithGroup(name string) slog.Handler {
	return h
}

func colorize(level slog.Level, msg string) string {
	switch level {
	case slog.LevelDebug:
		return "\033[37m" + msg + "\033[0m" // gray
	case slog.LevelInfo:
		return "\033[34m" + msg + "\033[0m" // blue
	case slog.LevelWarn:
		return "\033[33m" + msg + "\033[0m" // yellow
	default:
		return "\033[31m" + msg + "\033[0m" // red
	}
}

func isTerminal(f *os.File) bool {
	fi, _ := f.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w}

		logger.Info("incoming request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"size", rec.size,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

func InitLogger() *slog.Logger {
	handler := NewJSONHandler(slog.LevelInfo)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
