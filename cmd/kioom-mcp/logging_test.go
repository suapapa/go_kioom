package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		wantJSON  bool
		wantError bool
	}{
		{name: "default text", format: "", wantJSON: false},
		{name: "text alias", format: "text", wantJSON: false},
		{name: "pretty alias", format: "pretty", wantJSON: false},
		{name: "json", format: "json", wantJSON: true},
		{name: "unknown", format: "yaml", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger, err := setupLogger(tt.format, &buf)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("setupLogger() error = %v", err)
			}

			logger.Info("hello", "session", "abc")
			out := buf.String()
			if out == "" {
				t.Fatal("expected log output")
			}
			if tt.wantJSON {
				if !strings.HasPrefix(out, "{") {
					t.Fatalf("expected JSON output, got %q", out)
				}
			} else if strings.HasPrefix(out, "{") {
				t.Fatalf("expected text output, got JSON %q", out)
			}
		})
	}
}

func TestFatalUsesLogger(t *testing.T) {
	var buf bytes.Buffer
	logger, err := setupLogger("text", &buf)
	if err != nil {
		t.Fatalf("setupLogger() error = %v", err)
	}

	exitCode := -1
	oldExit := exitFunc
	exitFunc = func(code int) { exitCode = code }
	defer func() { exitFunc = oldExit }()

	fatal(logger, "boom", "reason", "test")
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Fatalf("expected fatal log in buffer, got %q", buf.String())
	}
}
