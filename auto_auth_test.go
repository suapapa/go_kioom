package kioom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAutoTokenIssuance(t *testing.T) {
	var tokenCalls int64
	var apiCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			atomic.AddInt64(&tokenCalls, 1)
			res := TokenResponse{
				ReturnCode: 0,
				ReturnMsg:  "success",
				Token:      "auto-issued-token",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(res)
			return
		}

		if r.URL.Path == "/api/dostk/acnt" {
			atomic.AddInt64(&apiCalls, 1)
			auth := r.Header.Get("authorization")
			if auth != "Bearer auto-issued-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			res := AccountResponse{
				AcctNo:     "123456",
				ReturnCode: 0,
				ReturnMsg:  "success",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(res)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithAutoToken())
	client.baseURL = server.URL

	// Make an API call without setting a token first
	res, err := client.GetAccountNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.AcctNo != "123456" {
		t.Errorf("expected account number 123456, got %s", res.AcctNo)
	}

	if atomic.LoadInt64(&tokenCalls) != 1 {
		t.Errorf("expected 1 token call, got %d", tokenCalls)
	}
	if atomic.LoadInt64(&apiCalls) != 1 {
		t.Errorf("expected 1 api call, got %d", apiCalls)
	}
}

func TestConcurrentAutoTokenIssuance(t *testing.T) {
	var tokenCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			atomic.AddInt64(&tokenCalls, 1)
			res := TokenResponse{
				ReturnCode: 0,
				ReturnMsg:  "success",
				Token:      "shared-token",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(res)
			return
		}

		if r.URL.Path == "/api/dostk/acnt" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(AccountResponse{AcctNo: "123456"})
			return
		}
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithAutoToken())
	client.baseURL = server.URL

	var wg sync.WaitGroup
	const numConcurrent = 10
	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.GetAccountNumber(context.Background())
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt64(&tokenCalls); got != 1 {
		t.Errorf("expected exactly 1 token call due to sharing/caching, got %d", got)
	}
}

func TestAutoTokenRefreshOnReturnCode3(t *testing.T) {
	var tokenCalls int64
	var apiCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			atomic.AddInt64(&tokenCalls, 1)
			res := TokenResponse{
				ReturnCode: 0,
				ReturnMsg:  "success",
				Token:      "fresh-token",
				ExpiresDt:  time.Now().In(kstLocation).Add(time.Hour).Format("20060102150405"),
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(res)
			return
		}

		if r.URL.Path == "/api/dostk/acnt" {
			apiCallNum := atomic.AddInt64(&apiCalls, 1)
			auth := r.Header.Get("authorization")

			if apiCallNum == 1 {
				if auth != "Bearer expired-token" {
					t.Errorf("expected first call to use expired-token, got %s", auth)
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(AccountResponse{
					ReturnCode: 3,
					ReturnMsg:  "8005:Token이 유효하지 않습니다",
				})
				return
			}

			if auth != "Bearer fresh-token" {
				t.Errorf("expected retried call to use fresh-token, got %s", auth)
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(AccountResponse{
				AcctNo:     "123456",
				ReturnCode: 0,
				ReturnMsg:  "success",
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithAutoToken())
	client.baseURL = server.URL
	client.SetToken("expired-token")

	res, err := client.GetAccountNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.AcctNo != "123456" {
		t.Errorf("expected account number 123456, got %s", res.AcctNo)
	}

	if got := atomic.LoadInt64(&tokenCalls); got != 1 {
		t.Errorf("expected 1 token call for refresh, got %d", got)
	}
	if got := atomic.LoadInt64(&apiCalls); got != 2 {
		t.Errorf("expected 2 api calls (1 failed + 1 retried), got %d", got)
	}
}

func TestIssueTokenSkipsExpiredCache(t *testing.T) {
	var tokenCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth2/token" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		atomic.AddInt64(&tokenCalls, 1)
		res := TokenResponse{
			ReturnCode: 0,
			ReturnMsg:  "success",
			Token:      "new-token",
			ExpiresDt:  time.Now().In(kstLocation).Add(time.Hour).Format("20060102150405"),
		}
		_ = json.NewEncoder(w).Encode(res)
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithAutoToken())
	client.baseURL = server.URL
	client.setTokenWithExpiry("old-token", time.Now().Add(-time.Minute))

	res, err := client.IssueToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Token != "new-token" {
		t.Errorf("expected new-token, got %s", res.Token)
	}
	if res.ReturnMsg == "success (cached)" {
		t.Fatal("expected expired token to be refreshed, got cached response")
	}
	if got := atomic.LoadInt64(&tokenCalls); got != 1 {
		t.Errorf("expected 1 token issue call, got %d", got)
	}
}

func TestIssueTokenReturnsValidCache(t *testing.T) {
	var tokenCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&tokenCalls, 1)
		_ = json.NewEncoder(w).Encode(TokenResponse{Token: "should-not-call"})
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithAutoToken())
	client.baseURL = server.URL
	client.setTokenWithExpiry("cached-token", time.Now().Add(time.Hour))

	res, err := client.IssueToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Token != "cached-token" {
		t.Errorf("expected cached-token, got %s", res.Token)
	}
	if res.ReturnMsg != "success (cached)" {
		t.Errorf("expected cached response, got %q", res.ReturnMsg)
	}
	if got := atomic.LoadInt64(&tokenCalls); got != 0 {
		t.Errorf("expected 0 token issue calls, got %d", got)
	}
}

func TestAutoTokenRefreshOn401(t *testing.T) {
	var tokenCalls int64
	var apiCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			atomic.AddInt64(&tokenCalls, 1)
			res := TokenResponse{
				ReturnCode: 0,
				ReturnMsg:  "success",
				Token:      "fresh-token",
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(res)
			return
		}

		if r.URL.Path == "/api/dostk/acnt" {
			apiCallNum := atomic.AddInt64(&apiCalls, 1)
			auth := r.Header.Get("authorization")

			// Check that request body is received correctly
			var req DepositRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.QryTp != "3" {
				t.Errorf("expected QryTp to be 3, got %q", req.QryTp)
			}

			if apiCallNum == 1 {
				// First call fails with 401 Unauthorized
				if auth != "Bearer expired-token" {
					t.Errorf("expected first call to use expired-token, got %s", auth)
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Second call should succeed with fresh-token
			if auth != "Bearer fresh-token" {
				t.Errorf("expected retried call to use fresh-token, got %s", auth)
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(DepositResponse{
				Deposit: "1000",
			})
			return
		}
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithAutoToken())
	client.baseURL = server.URL
	client.SetToken("expired-token")

	// Use GetDeposit which sends a request body
	res, err := client.GetDeposit(context.Background(), &DepositRequest{QryTp: "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Deposit != "1000" {
		t.Errorf("expected deposit 1000, got %s", res.Deposit)
	}

	if got := atomic.LoadInt64(&tokenCalls); got != 1 {
		t.Errorf("expected 1 token call for refresh, got %d", got)
	}
	if got := atomic.LoadInt64(&apiCalls); got != 2 {
		t.Errorf("expected 2 api calls (1 failed + 1 retried), got %d", got)
	}
}
