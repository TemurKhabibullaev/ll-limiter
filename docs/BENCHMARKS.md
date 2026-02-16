# BENCHMARKS

Run:

```bash
go test -bench=. -benchmem ./...

Results:
goos: linux
goarch: amd64
pkg: github.com/TemurKhabibullaev/ll-limiter/bench
cpu: AMD EPYC 7763 64-Core Processor                
BenchmarkTokenBucket_Allow_Cost1-2              12131535               100.0 ns/op             0 B/op          0 allocs/op
BenchmarkTokenBucket_Allow_ManyKeys-2            9987994               121.3 ns/op             0 B/op          0 allocs/op
BenchmarkTokenBucket_Allow_Cost10-2             12088305                98.18 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/TemurKhabibullaev/ll-limiter/bench   3.943s
?       github.com/TemurKhabibullaev/ll-limiter/cmd/limiter     [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/clock  [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/httpapi        [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/limiter        [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/metrics        [no test files]

-----------------
NOTES
* ns/op measures average time per decision.
* allocs/op indicates allocations on the hot path (lower is better for low-latency systems).
* Benchmarks include single-key and many-key scenarios to observe map/lock behavior.

---

SLIDING WINDOW:

@TemurKhabibullaev ➜ /workspaces/ll-limiter (main) $ go test -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: github.com/TemurKhabibullaev/ll-limiter/bench
cpu: AMD EPYC 7763 64-Core Processor                
BenchmarkTokenBucket_Allow_Cost1-2               6735381               154.4 ns/op             0 B/op          0 allocs/op
BenchmarkTokenBucket_Allow_ManyKeys-2            6645008               151.1 ns/op             0 B/op          0 allocs/op
BenchmarkTokenBucket_Allow_Cost10-2              8253910               148.8 ns/op             0 B/op          0 allocs/op
BenchmarkSlidingWindow_Allow_Cost1-2             6439879               189.8 ns/op             0 B/op          0 allocs/op
BenchmarkSlidingWindow_Allow_ManyKeys-2          8495972               131.3 ns/op             1 B/op          0 allocs/op
PASS
ok      github.com/TemurKhabibullaev/ll-limiter/bench   6.502s
?       github.com/TemurKhabibullaev/ll-limiter/cmd/limiter     [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/clock  [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/httpapi        [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/limiter        [no test files]
?       github.com/TemurKhabibullaev/ll-limiter/internal/metrics        [no test files]
