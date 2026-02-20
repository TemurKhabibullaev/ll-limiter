// internal/limiter/shards.go
package limiter

import (
	"hash/fnv"
	"sync"
)

// shard holds a portion of the keyspace with its own lock.
type shard[V any] struct {
	mu sync.Mutex
	m  map[string]V
}

// shardedMap is a fixed-count set of shards.
// Each shard has its own lock + map to reduce contention.
type shardedMap[V any] struct {
	shards []shard[V]
	mask   uint64 // if len(shards) is power of two, we can mask instead of mod
}

// newShardedMap constructs a sharded map with n shards.
// n must be >= 1. Prefer a power of two (64/128/256) for faster indexing.
func newShardedMap[V any](n int) shardedMap[V] {
	if n < 1 {
		n = 1
	}

	// Round up to next power of two for masking, if not already.
	pow2 := 1
	for pow2 < n {
		pow2 <<= 1
	}
	n = pow2

	s := make([]shard[V], n)
	for i := 0; i < n; i++ {
		s[i].m = make(map[string]V)
	}

	return shardedMap[V]{shards: s, mask: uint64(n - 1)}
}

// shardFor returns the shard for a given key.
func (sm shardedMap[V]) shardFor(key string) *shard[V] {
	idx := sm.index(key)
	return &sm.shards[idx]
}

// index hashes the key and picks a shard index.
func (sm shardedMap[V]) index(key string) int {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	sum := h.Sum64()

	// Because n is a power of two, mask is safe and fast.
	return int(sum & sm.mask)
}
