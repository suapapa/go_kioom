package kioom

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		appKey    string
		secretKey string
		opts      []Option
		expected  string
	}{
		{"Live Domain", "app1", "sec1", nil, LiveDomain},
		{"Mock Domain", "app2", "sec2", []Option{WithMockDomain()}, MockDomain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.appKey, tt.secretKey, tt.opts...)
			if client.BaseURL() != tt.expected {
				t.Errorf("expected BaseURL %q, got %q", tt.expected, client.BaseURL())
			}
			// appKey and secretKey are unexported, we can't check them directly
			// but we can trust the constructor for now or add internal getters if needed
		})
	}
}

func TestSetToken(t *testing.T) {
	client := NewClient("app", "sec", WithMockDomain())
	token := "sample-token"
	client.SetToken(token)
	if client.Token() != token {
		t.Errorf("expected Token %q, got %q", token, client.Token())
	}
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	// In the new implementation, baseURL is unexported.
	// We might need to add a way to set it for testing, or use a custom HTTP client.
	// Since I added WithHTTPClient, I can't easily override the base URL unless I change the code.
	// Actually, I can just use a custom Option for testing if I want to keep baseURL unexported.

	client := NewClient("app", "sec")
	client.baseURL = server.URL // This works since we are in the same package

	req, err := client.newRequest(context.Background(), http.MethodGet, "/test", "api123", nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	var res struct {
		Success bool `json:"success"`
	}

	if err := client.do(req, &res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success {
		t.Errorf("expected Success to be true")
	}
}
