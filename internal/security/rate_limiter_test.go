package security

import (
	"testing"
	"time"
)

func TestRateLimiterBlocksAfterFiveAttempts(t *testing.T) {
	limiter := NewRateLimiter()

	for i := 0; i < 5; i++ {
		if !limiter.Allow("127.0.0.1") {
			t.Fatalf("attempt %d should have been allowed", i+1)
		}
	}

	if limiter.Allow("127.0.0.1") {
		t.Fatal("sixth attempt should have been blocked")
	}
}

func TestRateLimiterResetsWindow(t *testing.T) {
	limiter := NewRateLimiter()
	ip := "127.0.0.1"

	for i := 0; i < 5; i++ {
		limiter.Allow(ip)
	}

	limiter.attempts[ip].reset = time.Now().Add(-time.Second)

	if !limiter.Allow(ip) {
		t.Fatal("expected limiter to allow request after reset")
	}

	if limiter.attempts[ip].count != 1 {
		t.Fatalf("expected reset count to be 1, got %d", limiter.attempts[ip].count)
	}
}
