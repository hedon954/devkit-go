// Package cacheevict provides some cache eviction policy algorithms.
package cacheevict

// Cache defines the interface for a cache.
type Cache interface {
	// Add adds a key-value pair to the cache.
	Add(string, any)
	// Get retrieves the value associated with the given key from the cache.
	Get(string) (any, bool)
}

type cacheItem struct {
	key   string
	value any
}

// Policy is a type for cache eviction policies.
type Policy string

const (
	FIFO Policy = "fifo"
	LRU  Policy = "lru"
	LFU  Policy = "lfu"
	ARC  Policy = "arc"
)

type builder struct {
	policy   Policy
	capacity int
}

// Builder returns a new builder for building a cache.
func Builder() *builder {
	return &builder{}
}

// Policy sets the policy of the cache.
func (b *builder) Policy(policy Policy) *builder {
	b.policy = policy
	return b
}

// Capacity sets the capacity of the cache.
func (b *builder) Capacity(capacity int) *builder {
	b.capacity = capacity
	return b
}

// Build builds a new cache with the given policy and capacity.
func (b *builder) Build() Cache {
	if b.policy == "" || b.capacity <= 0 {
		panic("unspecified policy or capacity")
	}

	return New(b.policy, b.capacity)
}

// New creates a new cache with the given policy and capacity.
func New(policy Policy, capacity int) Cache {
	switch policy {
	case FIFO:
		return NewFIFOCache(capacity)
	case LRU:
		return NewLRUCache(capacity)
	case LFU:
		return NewLFUCache(capacity)
	case ARC:
		return NewARCCache(capacity)
	default:
		panic("unsupported policy: " + policy)
	}
}
