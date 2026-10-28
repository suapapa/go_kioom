package kioom

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// defaultRateLimits holds the per-TR request rates Kiwoom publishes, in
// requests per second, for endpoints whose limit is tighter than the general
// allowance. The client spaces calls to these TRs evenly so a loop or a burst
// of goroutines cannot trip the server-side 유량제어 (rate limit).
//
// Kiwoom tightened kt00018 (계좌평가잔고내역요청) from 5 req/s to 3 req/s
// effective 2026-10-29 18:00 KST.
var defaultRateLimits = map[string]int{
	API_GetAccountBalance: 3,
}

// rateLimiter spaces calls so that at most one request starts per interval.
type rateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

func newRateLimiter(perSecond int) *rateLimiter {
	return &rateLimiter{interval: time.Second / time.Duration(perSecond)}
}

// wait blocks until the next request slot is available or ctx is done.
func (l *rateLimiter) wait(ctx context.Context) error {
	l.mu.Lock()
	now := time.Now()
	if l.next.Before(now) {
		l.next = now
	}
	delay := l.next.Sub(now)
	l.next = l.next.Add(l.interval)
	l.mu.Unlock()

	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// WithRateLimit overrides the client-side request rate for a TR code.
// perSecond is the maximum number of requests per second; values <= 0 disable
// throttling for that TR and leave the limit to the server.
func WithRateLimit(apiID string, perSecond int) Option {
	return func(c *Client) {
		c.rateLimits[apiID] = perSecond
	}
}

// WithoutRateLimit disables all client-side throttling. Use it only when the
// caller already paces requests itself; Kiwoom rejects calls above its limits.
func WithoutRateLimit() Option {
	return func(c *Client) {
		c.rateLimitDisabled = true
	}
}

// waitRateLimit blocks until a request for apiID is allowed by the client-side
// rate limit.
func (c *Client) waitRateLimit(ctx context.Context, apiID string) error {
	if apiID == "" || c.rateLimitDisabled {
		return nil
	}
	perSecond, ok := c.rateLimits[apiID]
	if !ok || perSecond <= 0 {
		return nil
	}

	c.rateMu.Lock()
	limiter, ok := c.rateLimiters[apiID]
	if !ok {
		limiter = newRateLimiter(perSecond)
		c.rateLimiters[apiID] = limiter
	}
	c.rateMu.Unlock()

	if err := limiter.wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait for %s: %w", apiID, err)
	}
	return nil
}
