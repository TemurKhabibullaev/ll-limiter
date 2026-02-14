ll-limiter

A low-latency token bucket rate limiter service written in Go.

Built as a systems and SRE-focused lab project to explore:

* concurrency design

* request-level observability

* low-latency HTTP services

* production-style metrics instrumentation

🚀 Features

Token Bucket rate limiting

Per-key isolation

HTTP API

Health endpoint (/healthz)

Prometheus metrics (/metrics)

Configurable via environment variables

Concurrency-safe implementation

📐 Architecture Overview
Client
   ↓
HTTP API (net/http)
   ↓
Token Bucket (in-memory)
   ↓
Decision (Allowed / Rejected)
   ↓
Prometheus Metrics


Single-process, in-memory design optimized for simplicity and low latency.

📡 API
GET /v1/allow

Query parameters:

Parameter	Required	Default
key	Yes	—
cost	No	1

Example:

curl "http://127.0.0.1:8080/v1/allow?key=user1&cost=1"


Response:

{
  "allowed": true,
  "remaining": 99,
  "retry_after_ms": 0,
  "algorithm": "token_bucket"
}

GET /healthz

Returns:

ok


Used for liveness checks.

GET /metrics

Exposes Prometheus metrics including:

request counters

allowed/rejected counters

latency histogram

in-flight requests

Go runtime metrics

⚙️ Configuration

Environment variables:

Variable	Default	Description
RATE_PER_SEC	50	Token refill rate
BURST	100	Maximum bucket size
PORT	8080	HTTP port

Example:

RATE_PER_SEC=100 BURST=200 make run

🧠 Design Principles

O(1) decision per request

Lock-protected token bucket

Lazy expiration for unused keys

No background cleanup goroutines

Deterministic behavior

Metrics-first design

📊 Observability

Example metrics:

ll_limiter_requests_total
ll_limiter_allowed_total
ll_limiter_rejected_total
ll_limiter_request_duration_seconds
ll_limiter_in_flight


Designed to be scrape-ready in real monitoring systems.

🧪 Local Development

Start server:

make run


Test:

curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/v1/allow?key=test
curl http://127.0.0.1:8080/metrics

🔮 Future Enhancements

Sliding window implementation

Redis-backed distributed limiter

Sharded bucket map to reduce lock contention

Benchmark suite

Dockerfile

Kubernetes deployment example

Save.

Commit:

git add README.md
git commit -m "add professional project documentation"
git push
