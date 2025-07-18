package utils

import (
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests int
	interval    time.Duration

	mu        sync.Mutex
	remaining int
	resetAt   time.Time
}

func NewRateLimiter(maxRequests int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		interval:    interval,
		remaining:   maxRequests,
		resetAt:     time.Now().Add(interval),
	}
}

func (rl *RateLimiter) Enqueue(job func()) {
	rl.acquire()
	go job()
}

func (rl *RateLimiter) acquire() {
	for {
		rl.mu.Lock()
		now := time.Now()

		if now.After(rl.resetAt) {
			rl.remaining = rl.maxRequests
			rl.resetAt = now.Add(rl.interval)
		}

		if rl.remaining > 0 {
			rl.remaining--
			rl.mu.Unlock()
			return
		}

		wait := rl.resetAt.Sub(now)
		rl.mu.Unlock()
		time.Sleep(wait)
	}
}
