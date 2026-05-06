package kioomenv

import (
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
