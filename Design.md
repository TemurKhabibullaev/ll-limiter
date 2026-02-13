# Design

## Goals
- Low-latency Allow() path
- Simple HTTP API
- Extensible algorithms/backends

## Current
- Token Bucket
- In-memory map + mutex
- Per-key buckets

## Next
- Prometheus metrics
- Benchmarks (p95/p99)
- Sharded map for reduced contention
- Optional Redis backend
