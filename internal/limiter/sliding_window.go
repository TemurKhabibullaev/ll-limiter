package limiter

import (
	"sync"
	"time"
)

// Sliding-window limiter using a timestamp log per key.
// It is O(k) per request where k = events in window for that key.
// Great for correctness + comparison
type SlidingWindow struct {
	mu     sync.Mutex
	clock  func() time.Time
	limit  int64
	window time.Duration
	keys   map[string][]time.Time
}

func NewSlidingWindow(nowFn func() time.Time, limit int64, window time.Duration) *SlidingWindow {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &SlidingWindow{
		clock:  nowFn,
		limit:  limit,
		window: window,
		keys:   make(map[string][]time.Time),
	}
}

func (s *SlidingWindow) Allow(key string, cost int64) Decision {
	if cost <= 0 {
		cost = 1
	}

	now := s.clock()
	cutoff := now.Add(-s.window)

	s.mu.Lock()
	defer s.mu.Unlock()

	events := s.keys[key]

	// Prune old events
	keepFrom := 0
	for keepFrom < len(events) && events[keepFrom].Before(cutoff) {
		keepFrom++
	}
	if keepFrom > 0 {
		events = events[keepFrom:]
	}

	used := int64(len(events))
	remaining := s.limit - used

	// Need "cost" capacity
	if remaining < cost {
		// RetryAfter = time until oldest event falls out of window
		var retry time.Duration
		var reset time.Duration
		if len(events) > 0 {
			oldest := events[0]
			resetAt := oldest.Add(s.window)
			if resetAt.After(now) {
				retry = resetAt.Sub(now)
				reset = retry
			}
		}
		// Save pruned slice back
		s.keys[key] = events

		return Decision{
			Allowed:    false,
			Remaining:  max64(remaining, 0),
			RetryAfter: retry,
			ResetAfter: reset,
			Limit:      s.limit,
			Burst:      s.limit,
			Algorithm:  "sliding_window",
		}
	}

	// Record "cost" events (simple model: 1 event = 1 cost unit)
	// For higher precision you can store cost-weighted entries later.
	for i := int64(0); i < cost; i++ {
		events = append(events, now)
	}

	s.keys[key] = events
	newUsed := int64(len(events))
	newRemaining := s.limit - newUsed

	return Decision{
		Allowed:    true,
		Remaining:  max64(newRemaining, 0),
		RetryAfter: 0,
		ResetAfter: s.window, // approximate: window horizon
		Limit:      s.limit,
		Burst:      s.limit,
		Algorithm:  "sliding_window",
	}
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
