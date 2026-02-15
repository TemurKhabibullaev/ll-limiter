# ll-limiter

Low-latency token bucket rate limiter service written in Go.

Built as a production systems imitation and SRE-focused lab project to explore:

- Concurrency design
- Low-latency HTTP services
- Observability and metrics
- Production-style API behavior

---

## 🚀 Features

- Token Bucket rate limiting
- Per-key isolation
- HTTP API
- Health endpoint (`/healthz`)
- Prometheus metrics (`/metrics`)
- Environment-based configuration
- Concurrency-safe implementation

---

## 📐 Architecture Overview

Client --> HTTP API (net/http) --> Token Bucket (in-memory) --> Decision (Allowed / Rejected) --> Prometheus Metrics


Single-process, in-memory architecture optimized for simplicity and low latency.

---

## 📡 API
### GET `/v1/allow`

Query parameters:

| Parameter | Required | Default |
|-----------|----------|----------|
| key       | Yes      | —        |
| cost      | No       | 1        |

Example:

```bash
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

| Variable     | Default | Description         |
| ------------ | ------- | ------------------- |
| RATE_PER_SEC | 50      | Token refill rate   |
| BURST        | 100     | Maximum bucket size |
| PORT         | 8080    | HTTP port           |

Example:

RATE_PER_SEC=100 BURST=200 make run

---

🧠 Design Principles
* O(1) decision per request
* Lock-protected token bucket
* Lazy expiration for unused keys
* No background cleanup goroutines
* Deterministic behavior
* Metrics-first design

---

📊 Observability
Example metrics:

ll_limiter_requests_total
ll_limiter_allowed_total
ll_limiter_rejected_total
ll_limiter_request_duration_seconds
ll_limiter_in_flight


Designed to be scrape-ready in real monitoring systems.

---

🧪 Local Development
Start server:

make run

BASE="http://127.0.0.1:8080"

# Health check
curl -i "$BASE/healthz"

Response:
HTTP/1.1 200 OK
Date: Sun, 15 Feb 2026 19:35:26 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

# Rate limit request
curl -i "$BASE/v1/allow?key=test"

Response:
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 15 Feb 2026 19:35:51 GMT
Content-Length: 78
{"allowed":true,"remaining":99,"retry_after_ms":0,"algorithm":"token_bucket"}

# Metrics
curl -i "$BASE/metrics" | head -n 20

Response:
% Total % Received % Xferd Average Speed Time Time Time Current
Dload Upload Total Spent Left Speed 0 0 0 0 0 0 0 0 --:--:-- --:--:-- --:--:--
0HTTP/1.1 200 OK Content-Type: text/plain; version=0.0.4; charset=utf-8; escaping=underscores Date: Sun, 15 Feb 2026 19:36:10 GMT Transfer-Encoding: chunked # HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles. # TYPE go_gc_duration_seconds summary go_gc_duration_seconds{quantile="0"} 0 go_gc_duration_seconds{quantile="0.25"} 0 go_gc_duration_seconds{quantile="0.5"} 0 go_gc_duration_seconds{quantile="0.75"} 0 go_gc_duration_seconds{quantile="1"} 0 go_gc_duration_seconds_sum 0 go_gc_duration_seconds_count 0 # HELP go_gc_gogc_percent Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent. # TYPE go_gc_gogc_percent gauge go_gc_gogc_percent 100 # HELP go_gc_gomemlimit_bytes Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variablehe runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes. 10# TYPE go_gc_gomemlimit_bytes gauge go_gc_gomemlimit_bytes 9.223372036854776e+18 0 10215 0 10215 0 0 4775k 0 --:--:-- --:--:-- --:--:-- 4987k


---

🔮 Future Enhancements
Sliding window implementation
Redis-backed distributed limiter
Sharded bucket map to reduce lock contention
Benchmark suite
Dockerfile
Kubernetes deployment example


git add README.md
git commit -m "add professional project documentation"
git push

---

## 📚 Documentation
- [Architecture](docs/ARCHITECTURE.md)
- [Design Decisions](docs/DESIGN.md)

