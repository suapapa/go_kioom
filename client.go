// Package kioom provides a Go wrapper for the Kiwoom Open API (REST).
// It handles authentication, stock information retrieval, and account management.
package kioom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// LiveDomain is the base URL for the live Kiwoom investment environment.
	LiveDomain = "https://api.kiwoom.com"
	// MockDomain is the base URL for the mock/simulated Kiwoom investment environment.
	MockDomain = "https://mockapi.kiwoom.com"
)

// Client is the entry point for interacting with the Kiwoom REST API.
// It keeps track of the base URL, authentication keys, and the issued token.
type Client struct {
	BaseURL    string       // API base URL (Live or Mock)
	AppKey     string       // Kiwoom App Key
	SecretKey  string       // Kiwoom Secret Key
	Token      string       // Issued access token (if available)
	HTTPClient *http.Client // Customizable HTTP client
}

// NewClient returns a new Kiwoom REST API client with the given credentials.
// If useMock is true, it uses the mock investment environment domain instead of the live API.
func NewClient(appKey, secretKey string, useMock bool) *Client {
	baseURL := LiveDomain
	if useMock {
		baseURL = MockDomain
	}
	return &Client{
		BaseURL:    baseURL,
		AppKey:     appKey,
		SecretKey:  secretKey,
		HTTPClient: &http.Client{},
	}
}

// SetToken sets the OAuth token manually. This is useful for reusing an existing,
// valid token across different client instances or application restarts.
func (c *Client) SetToken(token string) {
	c.Token = token
}

// newRequest constructs a generic http.Request with required headers for Kiwoom API
func (c *Client) newRequest(method, path, apiID string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.BaseURL+path, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if apiID != "" {
		req.Header.Set("api-id", apiID)
	}
	if c.Token != "" {
		req.Header.Set("authorization", "Bearer "+c.Token)
	}

	return req, nil
}

// do makes the actual HTTP request and parses into 'v'
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
