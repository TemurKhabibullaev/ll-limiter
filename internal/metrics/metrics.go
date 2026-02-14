package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	AllowedTotal    *prometheus.CounterVec
	RejectedTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	InFlight        prometheus.Gauge
}

func New() *Metrics {
	m := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "ll_limiter_requests_total", Help: "Total /v1/allow requests"},
			[]string{"algorithm"},
		),
		AllowedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "ll_limiter_allowed_total", Help: "Total allowed decisions"},
			[]string{"algorithm"},
		),
		RejectedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "ll_limiter_rejected_total", Help: "Total rejected decisions"},
			[]string{"algorithm"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: "ll_limiter_request_duration_seconds", Help: "Latency of /v1/allow"},
			[]string{"algorithm"},
		),
		InFlight: prometheus.NewGauge(prometheus.GaugeOpts{Name: "ll_limiter_in_flight", Help: "In-flight requests"}),
	}

	prometheus.MustRegister(m.RequestsTotal, m.AllowedTotal, m.RejectedTotal, m.RequestDuration, m.InFlight)
	return m
}
