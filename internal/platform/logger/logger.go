package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a structured JSON logger using Go's built-in slog package.
//
// [Archetype Convention Addition] — Justified by observability best practices
// and ISO/IEC 25010 Maintainability → Analysability.
// Structured logging complements the Prometheus/Grafana metrics stack
// specified in the SAD (ADR-10).
func New(level string) *slog.Logger {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	})

	return slog.New(handler)
}
