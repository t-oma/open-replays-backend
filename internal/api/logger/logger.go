// Package logger provides structured logging using slog.
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config holds logger configuration.
type Config struct {
	Level  string
	Format string
}

// New creates a new slog.Logger with a custom writer based on configuration.
func New(cfg Config, w io.Writer) (*slog.Logger, error) {
	if w == nil {
		w = os.Stdout
	}

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "text":
		handler = slog.NewTextHandler(w, opts)
	case "json", "":
		handler = slog.NewJSONHandler(w, opts)
	default:
		return nil, fmt.Errorf("unknown log format: %s", cfg.Format)
	}

	return slog.New(handler), nil
}

// parseLevel parses a string log level into slog.Level.
func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level: %s", s)
	}
}
