package limiter

import (
	"math"
	"sync"
	"time"

	"github.com/TemurKhabibullaev/ll-limiter/internal/clock"
)

type bucket struct {
	tokens float64
	last   time.Time
}

type TokenBucket struct {
	mu    sync.Mutex
	clk   clock.Clock
	rate  float64
	burst float64
	state map[string]*bucket
	ttl   time.Duration
}

func NewTokenBucket(clk clock.Clock, ratePerSec float64, burst int64, ttl time.Duration) *TokenBucket {
	return &TokenBucket{
		clk:   clk,
		rate:  ratePerSec,
		burst: float64(burst),
		state: make(map[string]*bucket),
		ttl:   ttl,
	}
}

func (tb *TokenBucket) Allow(key string, cost int64) Decision {
	if cost <= 0 {
		cost = 1
	}
	now := tb.clk.Now()

	tb.mu.Lock()
	defer tb.mu.Unlock()

	b := tb.state[key]
	if b == nil {
		b = &bucket{tokens: tb.burst, last: now}
		tb.state[key] = b
	}

	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens = math.Min(tb.burst, b.tokens+(elapsed*tb.rate))
		b.last = now
	}

	need := float64(cost)
	if b.tokens >= need {
		b.tokens -= need
		rem := int64(math.Floor(b.tokens))
		return Decision{
			Allowed:    true,
			Remaining:  rem,
			RetryAfter: 0,
			ResetAfter: time.Duration((tb.burst-b.tokens)/tb.rate) * time.Second,
			Limit:      int64(tb.rate),
			Burst:      int64(tb.burst),
			Algorithm:  "token_bucket",
		}
	}

	deficit := need - b.tokens
	retrySec := deficit / tb.rate
	return Decision{
		Allowed:    false,
		Remaining:  0,
		RetryAfter: time.Duration(retrySec * float64(time.Second)),
		ResetAfter: time.Duration((tb.burst-b.tokens)/tb.rate) * time.Second,
		Limit:      int64(tb.rate),
		Burst:      int64(tb.burst),
		Algorithm:  "token_bucket",
	}
}
