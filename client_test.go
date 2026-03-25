package kioom

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		appKey    string
		secretKey string
		useMock   bool
		expected  string
	}{
		{"Live Domain", "app1", "sec1", false, LiveDomain},
		{"Mock Domain", "app2", "sec2", true, MockDomain},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.appKey, tt.secretKey, tt.useMock)
			if client.BaseURL != tt.expected {
				t.Errorf("expected BaseURL %q, got %q", tt.expected, client.BaseURL)
			}
			if client.AppKey != tt.appKey {
				t.Errorf("expected AppKey %q, got %q", tt.appKey, client.AppKey)
			}
		})
	}
}

func TestSetToken(t *testing.T) {
	client := NewClient("app", "sec", true)
	token := "sample-token"
	client.SetToken(token)
	if client.Token != token {
		t.Errorf("expected Token %q, got %q", token, client.Token)
	}
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL

	req, err := client.newRequest(http.MethodGet, "/test", "api123", nil)
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
