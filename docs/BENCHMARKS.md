# BENCHMARKS

Run:

```bash
go test -bench=. -benchmem ./...

Baseline Results:

goos: linux
goarch: amd64
pkg: github.com/TemurKhabibullaev/ll-limiter/bench
cpu: AMD EPYC 7763 64-Core Processor

Token Bucket:

BenchmarkTokenBucket_Allow_Cost1-2        6735381   154.4 ns/op   0 B/op   0 allocs/op
BenchmarkTokenBucket_Allow_ManyKeys-2     6645008   151.1 ns/op   0 B/op   0 allocs/op
BenchmarkTokenBucket_Allow_Cost10-2       8253910   148.8 ns/op   0 B/op   0 allocs/op

Sliding Window:

BenchmarkSlidingWindow_Allow_Cost1-2      6439879   189.8 ns/op   0 B/op   0 allocs/op
BenchmarkSlidingWindow_Allow_ManyKeys-2   8495972   131.3 ns/op   1 B/op   0 allocs/op

Parallel (contention)

go test ./...

These benchmarks simulate concurrent traffic across ~1024 keys using `b.RunParallel`.

BenchmarkTokenBucket_Allow_Parallel-2           12422649                87.87 ns/op            0 B/op          0 allocs/op
BenchmarkSlidingWindow_Allow_Parallel-2          6948603               169.4 ns/op             1 B/op          0 allocs/op

Token Bucket is faster under parallel load in this build.
Both algorithms remain allocation-free in the hot path (0 allocs/op).
-----------------
NOTES
* ns/op measures average time per Allow() decision.
* allocs/op indicates allocations on the hot path (lower is better for low-latency systems, should stay near zero for low-latency hot paths).
* Benchmarks include single-key and many-key scenarios to observe map/lock behavior.
