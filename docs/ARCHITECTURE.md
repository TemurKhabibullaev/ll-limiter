ARCHITECTURE
Overview

ll-limiter is a single-process, in-memory HTTP rate limiting service built in Go.

The system is designed to be:

Low-latency

Deterministic

Concurrency-safe

Observability-first

It is intentionally minimal to isolate core rate-limiting behavior without external dependencies.

High-Level Flow
Client Request
      ↓
HTTP Server (net/http)
      ↓
HandleAllow()
      ↓
Token Bucket (per-key)
      ↓
Decision (Allow / Reject)
      ↓
JSON Response + Metrics

Components
1. HTTP Layer

Uses Go’s net/http standard library.

Responsibilities:

Parse query parameters

Call limiter logic

Return JSON response

Expose health endpoint

Expose Prometheus metrics

No third-party routing frameworks are used.

2. Token Bucket Engine

Core logic lives in the limiter package.

Each key maps to a bucket containing:

current token count

last refill timestamp

Refill strategy:

Tokens are replenished lazily on each request

No background refill goroutine

Refill calculated based on elapsed time

Time Complexity:

O(1) per request

3. Concurrency Model

Bucket storage protected by mutex

Single critical section per request

No lock sharding (intentionally simple baseline)

No background maintenance workers

Design tradeoff:

Simplicity over extreme scalability

Intended as a foundation for future sharded design

4. Observability Layer

Prometheus instrumentation is integrated at request boundary.

Metrics include:

Total requests

Allowed decisions

Rejected decisions

Request latency histogram

In-flight requests

Go runtime metrics

This allows:

Latency distribution analysis

Throughput monitoring

Capacity planning experiments

Networking Model

The server binds to:

0.0.0.0:<PORT>


This allows:

Local testing via 127.0.0.1

Codespaces forwarded port exposure

Container-based deployment compatibility

Memory Model

Buckets are stored in memory.

Key properties:

No persistence

State lost on restart

Suitable for experimentation or single-instance services

Future enhancement:

Redis or distributed backing store

Failure Modes

Possible failure scenarios:

Process restart clears all buckets

High key cardinality increases memory usage

Single mutex may become contention bottleneck under extreme load

These tradeoffs are intentional for clarity and iterative development.

Deployment Model

Currently:

Single binary

Single node

Stateless outside in-memory state

Suitable for:

Local development

Educational systems design

Lab environments

Containerized deployment

Design Philosophy

This project prioritizes:

Correctness

Transparency

Observability

Incremental complexity

Future versions may introduce:

Sharded maps

Multiple algorithms

Distributed coordination

Horizontal scalability