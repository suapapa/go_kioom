package kioom

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := AccountResponse{
			AcctNo:     "1234567890",
			ReturnCode: 0,
		}
		json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", true)
	client.BaseURL = server.URL

	res, err := client.GetAccountNumber()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.AcctNo != "1234567890" {
		t.Errorf("expected AcctNo 1234567890, got %s", res.AcctNo)
	}
}
