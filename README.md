# ll-limiter

[![ci](https://github.com/TemurKhabibullaev/ll-limiter/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/TemurKhabibullaev/ll-limiter/actions/workflows/ci.yml)
![Go Version](https://img.shields.io/badge/go-1.24-blue)
![Docker](https://img.shields.io/badge/docker-ready-blue)

Low-latency rate limiter service written in Go.

Supports multiple algorithms (Token Bucket and Sliding Window)
with runtime selection via environment configuration.

Built as a production-systems lab project focused on:

- Concurrency design
- Algorithm trade-offs
- Low-latency HTTP services
- Observability and metrics
- CI/CD & production-style delivery

---

## 🚀 Features

- Multiple algorithms:
  - Token Bucket
  - Sliding Window
- Runtime algorithm selection (env-based)
- Per-key isolation
- HTTP API
- Health endpoint (`/healthz`)
- Prometheus metrics (`/metrics`)
- Zero allocations in hot path (bench verified)
- Concurrency-safe implementation
- Benchmarks included
- GitHub Actions CI pipeline
- Dockerized production image (distroless, non-root)

---

## 📐 Architecture Overview
Algorithm selection occurs at startup:

ALGORITHM=token_bucket (default)
ALGORITHM=sliding_window

Client --> HTTP API (net/http) --> Selected Limiter (Token Bucket OR Sliding Window) --> Decision --> Prometheus Metrics


Single-process, in-memory architecture optimized for simplicity and low latency.

---

## 📡 API

### `GET /v1/allow`

**Query parameters:**

| Parameter | Required | Default |
|-----------|----------|----------|
| key       | Yes      | —        |
| cost      | No       | 1        |

**Example:**

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


## ⚙️ Configuration
Environment variables:

| Variable     | Default      | Description                         |
| ------------ | ------------ | ----------------------------------- |
| RATE_PER_SEC | 50           | Token refill rate / window rate     |
| BURST        | 100          | Maximum bucket size / window capacity|
| PORT         | 8080         | HTTP port                           |
| ALGORITHM    | token_bucket | Rate Limiting Algorithm             |
| WINDOW_MS    | 1000         | Sliding window size in ms (sliding window only) |

Example:

RATE_PER_SEC=100 BURST=200 make run

## Sliding Window example (10 requests per 1 second window):

```bash
ALGORITHM=sliding_window WINDOW_MS=1000 BURST=10 make run

ALGORITHM=sliding_window RATE_PER_SEC=10 BURST=10 make run

## 🧮 Algorithms

This service supports multiple rate-limiting algorithms via `ALGORITHM`.

| Algorithm | Best for | Behavior |
|-----------|----------|----------|
| `token_bucket` | smoothing bursts | tokens refill continuously; allows short bursts up to `BURST` |
| `sliding_window` | strict per-window enforcement | enforces max `BURST` events per `WINDOW_MS`; returns accurate `retry_after_ms` |

Quick sliding-window check:

BASE="http://127.0.0.1:8080"
for i in $(seq 1 15); do curl -s "$BASE/v1/allow?key=x" ; echo; done

Response:

{"allowed":true,"remaining":9,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":8,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":7,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":6,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":5,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":4,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":3,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":2,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":1,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":true,"remaining":0,"retry_after_ms":0,"algorithm":"sliding_window"}
{"allowed":false,"remaining":0,"retry_after_ms":911,"algorithm":"sliding_window"}
{"allowed":false,"remaining":0,"retry_after_ms":900,"algorithm":"sliding_window"}
{"allowed":false,"remaining":0,"retry_after_ms":889,"algorithm":"sliding_window"}
{"allowed":false,"remaining":0,"retry_after_ms":879,"algorithm":"sliding_window"}
{"allowed":false,"remaining":0,"retry_after_ms":869,"algorithm":"sliding_window"}

---
## 📈 Benchmarks (Linux amd64)

Token Bucket:
- ~150 ns/op
- 0 allocs/op

Sliding Window:
- ~130–190 ns/op
- 0 allocs/op

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

---

## 🐳 Run with Docker

### Build image

docker build -t ll-limiter:local .

Run Container

docker run --rm \
  -p 8080:8080 \
  -e RATE_PER_SEC=50 \
  -e BURST=100 \
  ll-limiter:local

Verify Endpoints

curl -i http://127.0.0.1:8080/healthz
curl -i "http://127.0.0.1:8080/v1/allow?key=docker"
curl -i http://127.0.0.1:8080/metrics | head -n 20

Docker image uses:

* Multi-stage build
* Distroless runtime image
* Non-root user
* Minimal attack surface

---

## 🔮 Roadmap

Planned enhancements:

- Sharded bucket map to reduce lock contention under high concurrency
- Redis-backed distributed limiter
- Kubernetes deployment example

---

## 📚 Documentation
- [Architecture](docs/ARCHITECTURE.md)
- [Design Decisions](docs/DESIGN.md)
- [Benchmarks](docs/BENCHMARKS.md)

