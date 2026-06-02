package kioomenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	getenv := func(k string) string {
		m := map[string]string{
			"KIOOM_APP_KEY":    " ak ",
			"KIOOM_SECRET_KEY": " sk ",
			"KIOOM_TOKEN":      " tok ",
			"KIOOM_MOCK":       " true ",
		}
		return m[k]
	}
	c := Load(getenv)
	if c.AppKey != "ak" || c.SecretKey != "sk" || c.Token != "tok" || !c.Mock {
		t.Fatalf("Load = %+v", c)
	}
}

func TestRequireAppKeys(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{name: "ok", cfg: Config{AppKey: "a", SecretKey: "b"}, wantErr: false},
		{name: "missing app", cfg: Config{SecretKey: "b"}, wantErr: true},
		{name: "missing secret", cfg: Config{AppKey: "a"}, wantErr: true},
		{name: "whitespace only", cfg: Config{AppKey: "  ", SecretKey: "b"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.RequireAppKeys()
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected %v", err)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), "required") {
				t.Fatalf("message: %v", err)
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	// We modify process env, so don't run in parallel with other tests that modify env.
	// However, TestLoad and TestRequireAppKeys don't use real env, so it's fine.

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	// Set an environment variable beforehand to test that it is NOT overridden
	os.Setenv("TEST_EXISTING_ENV", "original_val")
	defer os.Unsetenv("TEST_EXISTING_ENV")
	defer os.Unsetenv("TEST_NEW_ENV")
	defer os.Unsetenv("TEST_QUOTED_ENV")
	defer os.Unsetenv("TEST_SINGLE_QUOTED_ENV")

	content := `
# This is a comment
TEST_EXISTING_ENV = new_val
TEST_NEW_ENV=loaded_val
TEST_QUOTED_ENV="quoted_val"
TEST_SINGLE_QUOTED_ENV='single_quoted_val'
# Another comment
  
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test .env: %v", err)
	}

	if err := LoadEnvFile(envPath); err != nil {
		t.Fatalf("LoadEnvFile failed: %v", err)
	}

	// Verify existing env is not overwritten
	if got := os.Getenv("TEST_EXISTING_ENV"); got != "original_val" {
		t.Errorf("expected TEST_EXISTING_ENV = original_val, got %q", got)
	}

	// Verify new env is loaded
	if got := os.Getenv("TEST_NEW_ENV"); got != "loaded_val" {
		t.Errorf("expected TEST_NEW_ENV = loaded_val, got %q", got)
	}

	// Verify quoted env is loaded and stripped
	if got := os.Getenv("TEST_QUOTED_ENV"); got != "quoted_val" {
		t.Errorf("expected TEST_QUOTED_ENV = quoted_val, got %q", got)
	}

	// Verify single quoted env is loaded and stripped
	if got := os.Getenv("TEST_SINGLE_QUOTED_ENV"); got != "single_quoted_val" {
		t.Errorf("expected TEST_SINGLE_QUOTED_ENV = single_quoted_val, got %q", got)
	}

	// Verify loading non-existent file does not error
	if err := LoadEnvFile(filepath.Join(tmpDir, ".env.does.not.exist")); err != nil {
		t.Errorf("expected no error for non-existent file, got %v", err)
	}
}
