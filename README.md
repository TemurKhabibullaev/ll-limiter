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

---

## 🐳 Run with Docker

### Build image

```bash
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

## 🔮 Future Enhancements

Planned roadmap (ordered roughly by impact):

- **Second algorithm:** Sliding Window limiter (compare behavior vs token bucket under bursty traffic)
- **Concurrency scaling:** Sharded bucket map / lock striping to reduce contention under high parallelism
- **Distributed mode:** Redis-backed limiter backend (multi-instance coordination + consistent enforcement)
- **Kubernetes example:** Helm/Kustomize + Deployment/Service + Prometheus scrape annotations
- **Hardening:** TTL eviction for inactive keys, max-key guardrails, structured logging, request IDs

---

## 📚 Documentation
- [Architecture](docs/ARCHITECTURE.md)
- [Design Decisions](docs/DESIGN.md)
- [Benchmarks](docs/BENCHMARKS.md)

