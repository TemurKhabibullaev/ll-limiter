# DESIGN

## Purpose

`ll-limiter` implements a Token Bucket rate limiting algorithm optimized for:

- Predictable latency
- Deterministic request handling
- Minimal background complexity
- Clear observability boundaries

This document explains algorithmic decisions, concurrency tradeoffs, and system constraints.

---

## Why Token Bucket?

Several rate limiting algorithms exist:

- Fixed window
- Sliding window
- Leaky bucket
- Token bucket

Token Bucket was chosen because it:

1. Allows short bursts
2. Supports smooth refill over time
3. Provides constant-time decisions per request
4. Avoids maintaining per-request timestamp history

Compared to sliding windows:

- No need to store individual request timestamps
- No unbounded memory growth from event history
- No periodic cleanup job required

---

## Algorithm Model

Each key maps to a bucket with:

- `tokens` (float64)
- `lastRefillTime` (timestamp)

Parameters:

- `rate` → tokens per second
- `burst` → max token capacity

Decision flow per request:

1. Calculate elapsed time since `lastRefillTime`
2. Add `(elapsed * rate)` tokens
3. Cap tokens at `burst`
4. If `tokens >= cost` → allow and deduct `cost`
5. Else reject and return retry hint

Time complexity: **O(1)** per request.

---

## Refill Strategy

Refill is **lazy**, not background-driven:

- Tokens are recomputed only when a request arrives for a key
- No goroutines are scheduled for periodic refills

Advantages:

- No idle CPU usage
- No periodic wakeups
- Simpler lifecycle and fewer moving parts
- Better predictability in low traffic environments

Tradeoff:

- Buckets only update on access (acceptable for this project scope)

---

## Concurrency Design

Current baseline implementation:

- Buckets are stored in a shared in-memory map
- Access is protected with a mutex
- One critical section per request

This was chosen for:

- Maximum clarity
- Easy reasoning about correctness
- Establishing a measurable performance baseline

Expected behavior under load:

- Lock contention increases with concurrency
- Throughput will eventually be bounded by mutex contention

Future improvements:

- Sharded map (N partitions) to reduce contention
- Lock striping / per-key locking
- Alternative storage using `sync.Map` for specific access patterns

---

## Memory Characteristics

Memory usage scales with:

- Number of unique keys
- Size of bucket metadata per key

Current properties:

- In-memory only
- No persistence
- No eviction by default
- State resets on restart

Future improvements:

- TTL-based expiration for inactive keys
- Max-key guardrails
- LRU eviction
- External backing store (e.g., Redis) for distributed rate limiting

---

## Latency Characteristics

Request latency consists of:

- HTTP parsing / query param validation
- Mutex acquisition
- Token math + decision
- JSON encoding
- Prometheus instrumentation

Expected characteristics:

- Microsecond-scale compute cost per request
- Low allocation overhead
- Deterministic per-request work

Latency is observable via the histogram:

- `ll_limiter_request_duration_seconds`

---

## Observability Philosophy

Instrumentation is placed at the request boundary to capture:

- Throughput
- Decision outcomes
- Latency distribution
- Concurrency pressure

Custom metrics include:

- `ll_limiter_requests_total`
- `ll_limiter_allowed_total`
- `ll_limiter_rejected_total`
- `ll_limiter_request_duration_seconds`
- `ll_limiter_in_flight`

This enables:

- Burst behavior analysis
- Capacity experiments
- Load test validation
- Monitoring integrations in real stacks

---

## Known Limitations

This project is intentionally scoped as a single-process baseline.

Current limitations:

- Single instance only
- No distributed coordination
- No clock skew handling across nodes
- Single mutex bottleneck under high concurrency
- In-memory storage only

These tradeoffs are deliberate to keep:

- Behavior transparent
- Correctness easy to validate
- Iteration speed high

---

## Failure Modes

Possible degradation scenarios:

1. Very high key cardinality → memory growth
2. Heavy concurrency → mutex contention
3. GC pauses under allocation spikes
4. Process restart → state reset and burst re-allowance

---

## Future Design Directions

Planned explorations:

- Sliding window limiter (time-bucketed or ring-buffer based)
- Algorithm comparison suite (token bucket vs sliding window vs leaky bucket)
- Sharded limiter core
- Redis-backed distributed mode
- Adaptive rate limiting / dynamic token refill rates
- More advanced low-latency experiments

---

## Engineering Philosophy

This repo prioritizes:

1. Measurable behavior
2. Simplicity before scaling
3. Observability from day one
4. Incremental architectural evolution

The goal is not maximum scalability immediately, but a clean baseline that can be evolved deliberately and benchmarked at every step.
