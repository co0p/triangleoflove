package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

const maxViolations = 8 // caps backoff at 2^8 × window (≈ 4 days for a 1-minute window)

type ipState struct {
	count        int
	windowStart  time.Time
	violations   int
	blockedUntil time.Time
}

type rateLimiter struct {
	mu     sync.Mutex
	states map[string]*ipState
	limit  int
	window time.Duration
	now    func() time.Time
}

// RateLimit returns a middleware that enforces per-IP request throttling.
// limit is the maximum number of requests allowed within window.
// Violations beyond the limit trigger exponential backoff: the first excess
// request blocks the IP for window, the second for 2×window, the third for
// 4×window, and so on (capped at 2^maxViolations × window).
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	return RateLimitWithClock(limit, window, time.Now)
}

// RateLimitWithClock is like RateLimit but accepts a custom clock function.
// Use it in tests to control time without real sleeping.
func RateLimitWithClock(limit int, window time.Duration, clock func() time.Time) func(http.Handler) http.Handler {
	rl := &rateLimiter{
		states: make(map[string]*ipState),
		limit:  limit,
		window: window,
		now:    clock,
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.allow(clientIP(r)) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "too many requests"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()

	state, ok := rl.states[ip]
	if !ok {
		state = &ipState{windowStart: now}
		rl.states[ip] = state
	}

	// Reject immediately if inside an active backoff block.
	if now.Before(state.blockedUntil) {
		return false
	}

	// Reset the request counter when the window has elapsed.
	// Violations are intentionally not reset so backoff persists across windows.
	if now.After(state.windowStart.Add(rl.window)) {
		state.count = 0
		state.windowStart = now
	}

	state.count++
	if state.count > rl.limit {
		state.violations++
		if state.violations > maxViolations {
			state.violations = maxViolations
		}
		backoff := rl.window * time.Duration(1<<uint(state.violations-1))
		state.blockedUntil = now.Add(backoff)
		return false
	}

	return true
}

// clientIP returns the originating client IP from the request.
// It reads the leftmost value of X-Forwarded-For (the original client IP as
// set by Railway's edge proxy) and falls back to RemoteAddr when the header
// is absent.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For may be a comma-separated list: "client, proxy1, proxy2".
		// The leftmost entry is the original client IP.
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	// Strip the port from "host:port" or "[ipv6]:port".
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i != -1 {
		return addr[:i]
	}
	return addr
}
