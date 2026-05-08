package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"triangleoflove/backend/internal/web"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// Acceptance test — maps directly to the increment.md criterion.
func TestRateLimit_GivenExcessiveRequests_WhenRegisterOrLoginHit_ThenReturns429(t *testing.T) {
	const limit = 3
	handler := web.RateLimit(limit, time.Minute)(okHandler())

	makeReq := func(xff string) int {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/register", nil)
		r.Header.Set("X-Forwarded-For", xff)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Code
	}

	// First `limit` requests from the same IP must pass.
	for i := 0; i < limit; i++ {
		if got := makeReq("1.2.3.4"); got != http.StatusOK {
			t.Fatalf("request %d: want 200, got %d", i+1, got)
		}
	}

	// Request beyond the limit must be rejected.
	if got := makeReq("1.2.3.4"); got != http.StatusTooManyRequests {
		t.Fatalf("want 429, got %d", got)
	}
}

// Different IPs must not share quota.
func TestRateLimit_GivenDifferentIPs_WhenOneExceedsLimit_ThenOtherStillAllowed(t *testing.T) {
	handler := web.RateLimit(1, time.Minute)(okHandler())

	makeReq := func(xff string) int {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/register", nil)
		r.Header.Set("X-Forwarded-For", xff)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Code
	}

	makeReq("10.0.0.1") // consume quota for IP A
	if got := makeReq("10.0.0.1"); got != http.StatusTooManyRequests {
		t.Fatalf("IP A: want 429, got %d", got)
	}
	// IP B should still be allowed.
	if got := makeReq("10.0.0.2"); got != http.StatusOK {
		t.Fatalf("IP B: want 200, got %d", got)
	}
}

// When X-Forwarded-For is absent, RemoteAddr is used as the client IP.
func TestRateLimit_GivenNoXForwardedFor_WhenRemoteAddrUsed_ThenRateLimitedByRemoteAddr(t *testing.T) {
	handler := web.RateLimit(1, time.Minute)(okHandler())

	makeReq := func() int {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
		// httptest.NewRequest sets RemoteAddr to "192.0.2.1:1234" by default.
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Code
	}

	if got := makeReq(); got != http.StatusOK {
		t.Fatalf("first request: want 200, got %d", got)
	}
	if got := makeReq(); got != http.StatusTooManyRequests {
		t.Fatalf("second request: want 429, got %d", got)
	}
}

// When X-Forwarded-For contains multiple IPs (client, proxy…) the leftmost (client) IP governs.
func TestRateLimit_GivenXForwardedForWithMultipleIPs_WhenFirstIPDiffers_ThenBucketedByFirstIP(t *testing.T) {
	handler := web.RateLimit(1, time.Minute)(okHandler())

	makeReq := func(xff string) int {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/register", nil)
		r.Header.Set("X-Forwarded-For", xff)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Code
	}

	// IP A via shared proxy.
	makeReq("10.0.0.1, 172.16.0.1") // consume quota
	if got := makeReq("10.0.0.1, 172.16.0.1"); got != http.StatusTooManyRequests {
		t.Fatalf("same client IP: want 429, got %d", got)
	}
	// Different client IP through the same proxy — must be a separate bucket.
	if got := makeReq("10.0.0.2, 172.16.0.1"); got != http.StatusOK {
		t.Fatalf("different client IP: want 200, got %d", got)
	}
}

// Each successive violation doubles the block duration (exponential backoff).
// Uses RateLimitWithClock so time can be controlled without real sleeping.
func TestRateLimit_GivenRepeatedViolations_WhenBackoffDoubles_ThenSecondBlockIsLonger(t *testing.T) {
	limit := 2
	window := 10 * time.Second

	ts := time.Now()
	clock := func() time.Time { return ts }

	handler := web.RateLimitWithClock(limit, window, clock)(okHandler())

	makeReq := func() int {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/register", nil)
		r.Header.Set("X-Forwarded-For", "1.2.3.4")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w.Code
	}

	// First violation: consume the limit then trigger 429.
	makeReq() // count=1
	makeReq() // count=2 (at limit)
	if got := makeReq(); got != http.StatusTooManyRequests {
		t.Fatalf("1st violation: want 429, got %d", got)
	}
	// Blocked for 1×window (10s).

	// Advance past the first block (window + 1ms).
	ts = ts.Add(window + time.Millisecond)

	// Second violation within the new window.
	makeReq() // count=1
	makeReq() // count=2 (at limit)
	if got := makeReq(); got != http.StatusTooManyRequests {
		t.Fatalf("2nd violation: want 429, got %d", got)
	}
	// Blocked for 2×window (20s).

	// One window later (10s into the 20s block) — still blocked.
	ts = ts.Add(window)
	if got := makeReq(); got != http.StatusTooManyRequests {
		t.Fatalf("mid 2nd block: want 429, got %d", got)
	}

	// Advance past the full 2×window block.
	ts = ts.Add(window + time.Millisecond)
	if got := makeReq(); got != http.StatusOK {
		t.Fatalf("after 2nd block expires: want 200, got %d", got)
	}
}
