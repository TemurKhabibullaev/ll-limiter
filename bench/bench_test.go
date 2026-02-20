package bench

import (
	"github.com/TemurKhabibullaev/ll-limiter/internal/clock"
	"github.com/TemurKhabibullaev/ll-limiter/internal/limiter"
	"strconv"
	"testing"
	"time"
)

func BenchmarkTokenBucket_Allow_Parallel(b *testing.B) {
	tb := limiter.NewTokenBucket(clock.RealClock{}, 50, 100, 10*time.Minute)

	// Precompute keys to avoid allocations in the benchmark itself
	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = tb.Allow(keys[i&1023], 1)
			i++
		}
	})
}

func BenchmarkSlidingWindow_Allow_Parallel(b *testing.B) {
	sw := limiter.NewSlidingWindow(time.Now, 100, 1*time.Second)

	// Precompute keys to avoid allocations in the benchmark itself
	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = sw.Allow(keys[i&1023], 1)
			i++
		}
	})
}

func BenchmarkTokenBucket_Allow_Cost1(b *testing.B) {
	tb := limiter.NewTokenBucket(clock.RealClock{}, 50, 100, 10*time.Minute)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = tb.Allow("bench-key", 1)
	}
}

func BenchmarkTokenBucket_Allow_ManyKeys(b *testing.B) {
	tb := limiter.NewTokenBucket(clock.RealClock{}, 50, 100, 10*time.Minute)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "bench-key-" + strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = tb.Allow(keys[i%len(keys)], 1)
	}
}

func BenchmarkTokenBucket_Allow_Cost10(b *testing.B) {
	tb := limiter.NewTokenBucket(clock.RealClock{}, 50, 100, 10*time.Minute)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = tb.Allow("bench-key", 10)
	}
}

func BenchmarkSlidingWindow_Allow_Cost1(b *testing.B) {
	sw := limiter.NewSlidingWindow(time.Now, 100, 1*time.Second)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.Allow("bench-key", 1)
	}
}

func BenchmarkSlidingWindow_Allow_ManyKeys(b *testing.B) {
	sw := limiter.NewSlidingWindow(time.Now, 100, 1*time.Second)

	keys := make([]string, 1024)
	for i := range keys {
		keys[i] = "bench-key-" + strconv.Itoa(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sw.Allow(keys[i%len(keys)], 1)
	}
}
