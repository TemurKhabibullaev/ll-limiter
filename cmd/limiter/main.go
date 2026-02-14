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

	tb := limiter.NewTokenBucket(clock.RealClock{}, rate, burst, 10*time.Minute)
	m := metrics.New()
	srv := httpapi.Server{L: tb, M: m}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/allow", srv.HandleAllow)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	mux.Handle("/metrics", promhttp.Handler())

	addr := "0.0.0.0:" + strconv.FormatInt(port, 10)
	log.Printf("ll-limiter listening on %s (rate=%.2f burst=%d)", addr, rate, burst)
	log.Fatal(http.ListenAndServe(addr, mux))
}
