package limiter

import "time"

type Decision struct {
	Allowed      bool
	Remaining    int64
	RetryAfter   time.Duration
	ResetAfter   time.Duration
	Limit        int64
	Burst        int64
	Algorithm    string
}

type Limiter interface {
	Allow(key string, cost int64) Decision
}
