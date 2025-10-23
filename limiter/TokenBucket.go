package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Tokens int
	Rate int
	Capacity int
	lastRefillTime time.Time
	mu sync.Mutex
}

func NewTokenBucket(rate int, capacity int) *TokenBucket {
	return &TokenBucket{
		Tokens: capacity,
		Rate: rate,
		Capacity: capacity,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()

	tb.Tokens += int(elapsed) * tb.Rate

	if tb.Tokens > tb.Capacity {
		tb.Tokens = tb.Capacity
	}
	tb.lastRefillTime = now

	if tb.Tokens >= 1 {
		tb.Tokens -= 1
		return true
	}
	return false
}

var buckets = make(map[string]*TokenBucket)
var mu sync.Mutex

