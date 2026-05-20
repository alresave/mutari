package security

import (
	"sync"
	"time"
)

type attempt struct {
	count int
	reset time.Time
}

// RateLimiter tracks login attempts per IP address within a short rolling window.
type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*attempt
}

// NewRateLimiter creates an in-memory limiter for authentication endpoints.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string]*attempt),
	}
}

// Allow records one attempt for the IP and reports whether it is still under the limit.
func (r *RateLimiter) Allow(ip string) bool {

	r.mu.Lock()
	defer r.mu.Unlock()

	a, exists := r.attempts[ip]

	now := time.Now()

	if !exists || now.After(a.reset) {

		r.attempts[ip] = &attempt{
			count: 1,
			reset: now.Add(time.Minute),
		}

		return true
	}

	if a.count >= 5 {
		return false
	}

	a.count++

	return true
}
