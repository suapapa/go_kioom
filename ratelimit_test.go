package kioom

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func newBalanceTestServer(t *testing.T, hits *int32, mu *sync.Mutex) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-id") != API_GetAccountBalance {
			t.Errorf("expected api-id %s, got %s", API_GetAccountBalance, r.Header.Get("api-id"))
		}
		if mu != nil && hits != nil {
			mu.Lock()
			*hits++
			mu.Unlock()
		}
		json.NewEncoder(w).Encode(AccountBalanceResponse{ReturnCode: 0, ReturnMsg: "ok"})
	}))
}

// TestAccountBalanceRateLimitDefault verifies kt00018 is throttled to 3 req/s
// (Kiwoom tightened the limit from 5 to 3 req/s on 2026-10-29 18:00 KST).
func TestAccountBalanceRateLimitDefault(t *testing.T) {
	server := newBalanceTestServer(t, nil, nil)
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	const calls = 4
	start := time.Now()
	for i := 0; i < calls; i++ {
		if _, err := client.GetAccountBalance(context.Background(), &AccountBalanceRequest{QryTp: "2"}); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
	elapsed := time.Since(start)

	// 4 calls at 3 req/s need at least 3 intervals ≈ 1s. Allow slack for timers.
	if min := 900 * time.Millisecond; elapsed < min {
		t.Errorf("4 calls completed in %v, expected at least %v (3 req/s spacing)", elapsed, min)
	}
}

// TestAccountBalanceRateLimitConcurrent verifies that concurrent callers share
// one limiter instead of each bursting independently.
func TestAccountBalanceRateLimitConcurrent(t *testing.T) {
	var (
		mu   sync.Mutex
		hits int32
	)
	server := newBalanceTestServer(t, &hits, &mu)
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain(), WithRateLimit(API_GetAccountBalance, 10))
	client.baseURL = server.URL

	const calls = 5
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < calls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.GetAccountBalance(context.Background(), &AccountBalanceRequest{QryTp: "2"}); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	// 5 calls at 10 req/s need at least 4 intervals = 400ms.
	if min := 350 * time.Millisecond; elapsed < min {
		t.Errorf("%d concurrent calls completed in %v, expected at least %v", calls, elapsed, min)
	}

	mu.Lock()
	got := hits
	mu.Unlock()
	if got != calls {
		t.Errorf("expected %d requests, got %d", calls, got)
	}
}

// TestAccountBalanceRateLimitDisabled verifies the opt-out leaves requests
// unthrottled for callers that pace themselves.
func TestAccountBalanceRateLimitDisabled(t *testing.T) {
	server := newBalanceTestServer(t, nil, nil)
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain(), WithoutRateLimit())
	client.baseURL = server.URL

	start := time.Now()
	for i := 0; i < 4; i++ {
		if _, err := client.GetAccountBalance(context.Background(), &AccountBalanceRequest{QryTp: "2"}); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("expected unthrottled calls, took %v", elapsed)
	}
}

// TestAccountBalanceRateLimitContextCanceled verifies a waiting call aborts
// with the context error instead of blocking past its deadline.
func TestAccountBalanceRateLimitContextCanceled(t *testing.T) {
	server := newBalanceTestServer(t, nil, nil)
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain(), WithRateLimit(API_GetAccountBalance, 1))
	client.baseURL = server.URL

	if _, err := client.GetAccountBalance(context.Background(), &AccountBalanceRequest{QryTp: "2"}); err != nil {
		t.Fatalf("first call: unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := client.GetAccountBalance(ctx, &AccountBalanceRequest{QryTp: "2"})
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > 700*time.Millisecond {
		t.Errorf("canceled wait took %v, expected it to return promptly", elapsed)
	}
}

// TestOtherTRsNotThrottled verifies only TRs with a documented limit are paced.
func TestOtherTRsNotThrottled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(AccountResponse{ReturnCode: 0})
	}))
	defer server.Close()

	client := NewClient("app", "sec", WithMockDomain())
	client.baseURL = server.URL

	start := time.Now()
	for i := 0; i < 4; i++ {
		if _, err := client.GetAccountNumber(context.Background()); err != nil {
			t.Fatalf("call %d: unexpected error: %v", i, err)
		}
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("unrelated TR appears throttled, took %v", elapsed)
	}
}
