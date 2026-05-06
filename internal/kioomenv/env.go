// Package kioomenv loads Kiwoom credential-related settings from environment variables,
// shared by cmd/kioom-cli, cmd/kioom-mcp, and cmd/examples.
package kioomenv

import (
	"errors"
	"strings"
)

// Known environment variables:
//   - KIOOM_APP_KEY
//   - KIOOM_SECRET_KEY
//   - KIOOM_TOKEN (optional bearer token)
//   - KIOOM_MOCK: when "true", callers typically use the mock API domain
//
// ErrMissingCredentials is returned when app or secret key is absent after Load.
var ErrMissingCredentials = errors.New("KIOOM_APP_KEY and KIOOM_SECRET_KEY are required")

// Config holds credential fields read from the process environment.
type Config struct {
	AppKey    string
	SecretKey string
	Token     string
	Mock      bool
}

// Load reads KIOOM_* variables using getenv (normally [os.Getenv]).
func Load(getenv func(string) string) Config {
	mock := strings.EqualFold(strings.TrimSpace(getenv("KIOOM_MOCK")), "true")
	return Config{
		AppKey:    strings.TrimSpace(getenv("KIOOM_APP_KEY")),
		SecretKey: strings.TrimSpace(getenv("KIOOM_SECRET_KEY")),
		Token:     strings.TrimSpace(getenv("KIOOM_TOKEN")),
		Mock:      mock,
	}
}

// RequireAppKeys returns [ErrMissingCredentials] if AppKey or SecretKey is empty.
func (c Config) RequireAppKeys() error {
	if strings.TrimSpace(c.AppKey) == "" || strings.TrimSpace(c.SecretKey) == "" {
		return ErrMissingCredentials
	}
	return nil
}
