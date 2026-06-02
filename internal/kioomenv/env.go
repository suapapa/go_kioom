// Package kioomenv loads Kiwoom credential-related settings from environment variables,
// shared by cmd/kioom-cli, cmd/kioom-mcp, and cmd/examples.
package kioomenv

import (
	"bufio"
	"errors"
	"os"
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

// LoadEnvFile reads a .env file from the specified path, parses KEY=VALUE pairs,
// and sets them in the environment using os.Setenv if they are not already set.
// If the file does not exist, it does nothing and returns nil.
func LoadEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Strip quotes if present
		if len(val) >= 2 {
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
				(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				val = val[1 : len(val)-1]
			}
		}

		// Only set if not already set in the environment to allow env overrides
		if os.Getenv(key) == "" {
			if err := os.Setenv(key, val); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
