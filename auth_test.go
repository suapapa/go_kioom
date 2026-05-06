package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIssueToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/oauth2/token" {
			t.Errorf("expected path /oauth2/token, got %s", r.URL.Path)
		}

		res := TokenResponse{
			ReturnCode: 0,
			ReturnMsg:  "success",
			Token:      "new-auth-token",
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	res, err := client.IssueToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Token != "new-auth-token" {
		t.Errorf("expected token new-auth-token, got %s", res.Token)
	}

	if client.Token() != "new-auth-token" {
		t.Errorf("expected client token to be updated, got %s", client.Token())
	}
}

func TestRevokeToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := RevokeResponse{
			ReturnCode: 0,
			ReturnMsg:  "success",
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL
	client.SetToken("old-token")

	res, err := client.RevokeToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ReturnCode != 0 {
		t.Errorf("expected ReturnCode 0, got %d", res.ReturnCode)
	}

	if client.Token() != "" {
		t.Errorf("expected client token to be cleared, got %s", client.Token())
	}
}
