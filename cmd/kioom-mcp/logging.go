package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const componentName = "kioom-mcp"

func setupLogger(format string, w io.Writer) (*slog.Logger, error) {
	if w == nil {
		w = os.Stderr
	}

	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var handler slog.Handler

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "text", "pretty":
		handler = slog.NewTextHandler(w, opts)
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	default:
		return nil, fmt.Errorf("unknown -log-format %q: use text or json", format)
	}

	logger := slog.New(handler).With("component", componentName)
	return logger, nil
}

var exitFunc = os.Exit

func fatal(logger *slog.Logger, msg string, args ...any) {
	logger.Error(msg, args...)
	exitFunc(1)
}
