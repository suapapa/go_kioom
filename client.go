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
	"sync"
	"time"
)

// Client is the entry point for interacting with the Kiwoom REST API.
type Client struct {
	baseURL    string
	appKey     string
	secretKey  string
	httpClient *http.Client

	mu    sync.RWMutex
	token string
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
func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

// Token returns the current access token.
func (c *Client) Token() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

// BaseURL returns the base URL of the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// newRequest constructs a generic http.Request with required headers for Kiwoom API.
func (c *Client) newRequest(ctx context.Context, method, path, apiID string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encode body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
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
