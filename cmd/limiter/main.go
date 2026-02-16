package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/TemurKhabibullaev/ll-limiter/internal/clock"
	"github.com/TemurKhabibullaev/ll-limiter/internal/httpapi"
	"github.com/TemurKhabibullaev/ll-limiter/internal/limiter"
	"github.com/TemurKhabibullaev/ll-limiter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getenvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

func getenvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

func main() {
	rate := getenvFloat("RATE_PER_SEC", 50)
	burst := getenvInt64("BURST", 100)
	port := getenvInt64("PORT", 8080)

	algo := os.Getenv("ALGORITHM")
	if algo == "" {
		algo = "token_bucket"
	}

	var L limiter.Limiter

	switch algo {
	case "sliding_window":
		windowMs := getenvInt64("WINDOW_MS", 1000)
		window := time.Duration(windowMs) * time.Millisecond
		// Sliding window: allow up to BURST events per window.
		L = limiter.NewSlidingWindow(time.Now, burst, window)
	case "token_bucket":
		fallthrough
	default:
		L = limiter.NewTokenBucket(clock.RealClock{}, rate, burst, 10*time.Minute)
	}

	m := metrics.New()
	srv := httpapi.Server{L: L, M: m}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/allow", srv.HandleAllow)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	mux.Handle("/metrics", promhttp.Handler())

	addr := "0.0.0.0:" + strconv.FormatInt(port, 10)
	log.Printf("ll-limiter listening on %s (rate=%.2f burst=%d)", addr, rate, burst)
	log.Fatal(http.ListenAndServe(addr, mux))
}
