// Package kioom provides a Go wrapper for the Kiwoom Open API (REST).
// It handles authentication, stock information retrieval, and account management.
package kioom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Client is the entry point for interacting with the Kiwoom REST API.
type Client struct {
	baseURL    string
	appKey     string
	secretKey  string
	httpClient *http.Client

	mu      sync.RWMutex
	token   string
	tokenMu sync.Mutex

	autoToken bool
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithMockDomain uses the mock investment environment domain.
func WithMockDomain() Option {
	return func(c *Client) {
		c.baseURL = MockDomain
	}
}

// WithAutoToken enables automatic token issuance and renewal.
func WithAutoToken() Option {
	return func(c *Client) {
		c.autoToken = true
	}
}

// NewClient returns a new Kiwoom REST API client with the given credentials.
func NewClient(appKey, secretKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:   LiveDomain,
		appKey:    appKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetToken sets the OAuth token manually.
// Also updates the environment variable KIOOM_TOKEN.
func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
	os.Setenv("KIOOM_TOKEN", token)
}

// Token returns the current access token.
// Checks environment variable KIOOM_TOKEN if internal token is empty.
func (c *Client) Token() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.token != "" {
		return c.token
	}
	return os.Getenv("KIOOM_TOKEN")
}

// BaseURL returns the base URL of the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// newRequest constructs a generic http.Request with required headers for Kiwoom API.
func (c *Client) newRequest(ctx context.Context, method, path, apiID string, body interface{}) (*http.Request, error) {
	// Avoid recursion/loop when calling token issue or revoke endpoints
	if c.autoToken && path != "/oauth2/token" && path != "/oauth2/revoke" {
		if c.Token() == "" {
			if _, err := c.IssueToken(ctx); err != nil {
				return nil, fmt.Errorf("auto-issue token: %w", err)
			}
		}
	}

	var bodyReader io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encode body: %w", err)
		}
		bodyReader = bytes.NewReader(buf.Bytes())
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if apiID != "" {
		req.Header.Set("api-id", apiID)
	}

	token := c.Token()
	if token != "" {
		req.Header.Set("authorization", "Bearer "+token)
	}

	return req, nil
}

// do makes the actual HTTP request and parses into 'v'.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if c.autoToken && resp.StatusCode == http.StatusUnauthorized && req.URL.Path != "/oauth2/token" && req.URL.Path != "/oauth2/revoke" {
		oldToken := req.Header.Get("authorization")
		if strings.HasPrefix(oldToken, "Bearer ") {
			oldToken = oldToken[7:]
		}

		c.mu.Lock()
		if c.token == oldToken {
			c.token = ""
		}
		c.mu.Unlock()

		_, err := c.IssueToken(req.Context())
		if err != nil {
			return fmt.Errorf("auto-refresh token failed: %w", err)
		}

		req.Header.Set("authorization", "Bearer "+c.Token())

		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return fmt.Errorf("get body for retry: %w", err)
			}
			req.Body = body
		}

		resp2, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("retry request: %w", err)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode < 200 || resp2.StatusCode >= 300 {
			return fmt.Errorf("API request failed with status: %s (retry)", resp2.Status)
		}

		if v != nil {
			if err := json.NewDecoder(resp2.Body).Decode(v); err != nil {
				return fmt.Errorf("decode response (retry): %w", err)
			}
		}
		return nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
