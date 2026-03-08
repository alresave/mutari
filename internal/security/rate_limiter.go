package security

import (
	"sync"
	"time"
)

type attempt struct {
	count int
	reset time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*attempt
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string]*attempt),
	}
}

func (r *RateLimiter) Allow(ip string) bool {

	r.mu.Lock()
	defer r.mu.Unlock()

	a, exists := r.attempts[ip]

	if !exists || time.Now().After(a.reset) {

		r.attempts[ip] = &attempt{
			count: 1,
			reset: time.Now().Add(time.Minute),
		}

		return true
	}

	if a.count >= 5 {
		return false
	}

	a.count++

	return true
}
