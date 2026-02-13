package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/TemurKhabibullaev/ll-limiter/internal/limiter"
)

type Server struct {
	L limiter.Limiter
}

type resp struct {
	Allowed    bool   `json:"allowed"`
	Remaining  int64  `json:"remaining"`
	RetryAfter int64  `json:"retry_after_ms"`
	Algorithm  string `json:"algorithm"`
}

func (s Server) HandleAllow(w http.ResponseWriter, r *http.Request) {
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
