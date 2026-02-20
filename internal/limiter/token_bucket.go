package limiter

import (
	"github.com/TemurKhabibullaev/ll-limiter/internal/clock"
	"math"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}

type TokenBucket struct {
	clk   clock.Clock
	rate  float64
	burst float64
	state shardedMap[*bucket]
	ttl   time.Duration
}

func NewTokenBucket(clk clock.Clock, ratePerSec float64, burst int64, ttl time.Duration) *TokenBucket {
	return &TokenBucket{
		clk:   clk,
		rate:  ratePerSec,
		burst: float64(burst),
		state: newShardedMap[*bucket](128),
		ttl:   ttl,
	}
}

func (tb *TokenBucket) Allow(key string, cost int64) Decision {
	if cost <= 0 {
		cost = 1
	}
	now := tb.clk.Now()
	sh := tb.state.shardFor(key)
	sh.mu.Lock()
	defer sh.mu.Unlock()

	b := sh.m[key]
	if b == nil {
		b = &bucket{tokens: tb.burst, last: now}
		sh.m[key] = b
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
