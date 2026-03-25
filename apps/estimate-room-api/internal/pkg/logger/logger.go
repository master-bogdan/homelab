// Package logger internal/logger/logger.go
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

var (
	defaultLogger *slog.Logger
	initOnce      sync.Once
)

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

	// Inject request_id if present
	if rid := GetRequestID(ctx); rid != "" {
		logEntry["request_id"] = rid
	}

	// collect attributes
	for _, a := range h.attrs {
		logEntry[a.Key] = normalizeLogValue(a.Value.Any())
	}
	r.Attrs(func(a slog.Attr) bool {
		logEntry[a.Key] = normalizeLogValue(a.Value.Any())
		return true
	})

	// JSON encode
	data, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	// add colors for local console
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

// --- helpers ---

func colorize(level slog.Level, msg string) string {
	switch level {
	case slog.LevelDebug:
		return "\033[37m" + msg + "\033[0m" // gray
	case slog.LevelInfo:
		return "\033[36m" + msg + "\033[0m" // cyan
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

func normalizeLogValue(value any) any {
	switch v := value.(type) {
	case nil:
		return nil
	case error:
		return v.Error()
	}

	if _, err := json.Marshal(value); err == nil {
		return value
	}

	return fmt.Sprint(value)
}

func Prefix(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	var prefix strings.Builder
	for _, part := range parts[:len(parts)-1] {
		trimmed := strings.ToUpper(strings.TrimSpace(part))
		if trimmed == "" {
			continue
		}
		prefix.WriteString("[")
		prefix.WriteString(trimmed)
		prefix.WriteString("]")
	}

	message := strings.TrimSpace(parts[len(parts)-1])
	if prefix.Len() == 0 {
		return message
	}
	if message == "" {
		return prefix.String()
	}

	return prefix.String() + " " + message
}

func InitLogger() *slog.Logger {
	initOnce.Do(func() {
		handler := NewJSONHandler(slog.LevelInfo)
		defaultLogger = slog.New(handler)
		slog.SetDefault(defaultLogger)
	})
	return defaultLogger
}

func L() *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	return InitLogger()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil || requestID == "" {
		return ctx
	}
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if rid, ok := ctx.Value(requestIDKey).(string); ok {
		return rid
	}
	return ""
}
