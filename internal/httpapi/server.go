package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/TemurKhabibullaev/ll-limiter/internal/limiter"
	"github.com/TemurKhabibullaev/ll-limiter/internal/metrics"
)

type Server struct {
	L limiter.Limiter
	M *metrics.Metrics
}

type resp struct {
	Allowed    bool   `json:"allowed"`
	Remaining  int64  `json:"remaining"`
	RetryAfter int64  `json:"retry_after_ms"`
	Algorithm  string `json:"algorithm"`
}

func (s Server) HandleAllow(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.M.InFlight.Inc()
	defer s.M.InFlight.Dec()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	cost := int64(1)
	if v := r.URL.Query().Get("cost"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			cost = n
		}
	}

	d := s.L.Allow(key, cost)
	alg := d.Algorithm
	s.M.RequestsTotal.WithLabelValues(alg).Inc()
	s.M.RequestDuration.WithLabelValues(alg).Observe(time.Since(start).Seconds())
	if d.Allowed {
		s.M.AllowedTotal.WithLabelValues(alg).Inc()
	} else {
		s.M.RejectedTotal.WithLabelValues(alg).Inc()
	}

	if !d.Allowed && d.RetryAfter > 0 {
		w.Header().Set("Retry-After", strconv.FormatInt(int64(d.RetryAfter.Seconds()), 10))
	}
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(resp{
		Allowed:    d.Allowed,
		Remaining:  d.Remaining,
		RetryAfter: int64(d.RetryAfter / time.Millisecond),
		Algorithm:  d.Algorithm,
	})
}
