package cacheevict

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache_New(t *testing.T) {
	t.Run("invalid capacity should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLRUCache(0)
		}, "Expected panic for invalid capacity")
	})

	t.Run("valid capacity should work", func(t *testing.T) {
		cache := NewLRUCache(3)
		assert.NotNil(t, cache, "Expected cache to be created")
	})
}

func TestLRUCache_AddAndGet(t *testing.T) {
	cache := NewLRUCache(3)

	// Add some items
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)

	// Retrieve them and ensure values are correct
	value, found := cache.Get("a")
	assert.True(t, found, "Expected to find key 'a'")
	assert.Equal(t, 1, value, "Expected value 1 for key 'a'")

	value, found = cache.Get("b")
	assert.True(t, found, "Expected to find key 'b'")
	assert.Equal(t, 2, value, "Expected value 2 for key 'b'")

	value, found = cache.Get("c")
	assert.True(t, found, "Expected to find key 'c'")
	assert.Equal(t, 3, value, "Expected value 3 for key 'c'")
}

func TestLRUCache_ExceedCapacity(t *testing.T) {
	cache := NewLRUCache(3)

	// Add more items than capacity allows
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)
	cache.Add("d", 4) // This should evict "a"

	// "a" should no longer be present, others should be present
	_, found := cache.Get("a")
	assert.False(t, found, "Expected 'a' to be evicted")

	value, found := cache.Get("b")
	assert.True(t, found, "Expected to find key 'b'")
	assert.Equal(t, 2, value, "Expected value 2 for key 'b'")

	value, found = cache.Get("c")
	assert.True(t, found, "Expected to find key 'c'")
	assert.Equal(t, 3, value, "Expected value 3 for key 'c'")

	value, found = cache.Get("d")
	assert.True(t, found, "Expected to find key 'd'")
	assert.Equal(t, 4, value, "Expected value 4 for key 'd'")
}

func TestLRUCache_RefreshExistingKey(t *testing.T) {
	cache := NewLRUCache(3)

	// Add some items
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)

	// Access key "a", making it the most recently used
	cache.Get("a")

	// Add another item, "b" should be evicted now, since "a" was accessed recently
	cache.Add("d", 4)

	_, found := cache.Get("b")
	assert.False(t, found, "Expected 'b' to be evicted")

	// "a", "c", and "d" should still be present
	value, found := cache.Get("a")
	assert.True(t, found, "Expected to find key 'a'")
	assert.Equal(t, 1, value, "Expected value 1 for key 'a'")

	value, found = cache.Get("c")
	assert.True(t, found, "Expected to find key 'c'")
	assert.Equal(t, 3, value, "Expected value 3 for key 'c'")

	value, found = cache.Get("d")
	assert.True(t, found, "Expected to find key 'd'")
	assert.Equal(t, 4, value, "Expected value 4 for key 'd'")
}

func TestLRUCache_UpdateExistingKey(t *testing.T) {
	cache := NewLRUCache(2)

	// Add some items
	cache.Add("a", 1)
	cache.Add("b", 2)

	// Update the value of "a"
	cache.Add("a", 10)

	// "a" should have the updated value, "b" should still be present
	value, found := cache.Get("a")
	assert.True(t, found, "Expected to find key 'a'")
	assert.Equal(t, 10, value, "Expected updated value 10 for key 'a'")

	value, found = cache.Get("b")
	assert.True(t, found, "Expected to find key 'b'")
	assert.Equal(t, 2, value, "Expected value 2 for key 'b'")
}

func TestLRUCache_EvictLeastRecentlyUsed(t *testing.T) {
	cache := NewLRUCache(2)

	// Add two items
	cache.Add("a", 1)
	cache.Add("b", 2)

	// Access key "a" so that "b" becomes the least recently used
	cache.Get("a")

	// Add a new item, "b" should be evicted
	cache.Add("c", 3)

	_, found := cache.Get("b")
	assert.False(t, found, "Expected 'b' to be evicted")

	// "a" and "c" should still be present
	value, found := cache.Get("a")
	assert.True(t, found, "Expected to find key 'a'")
	assert.Equal(t, 1, value, "Expected value 1 for key 'a'")

	value, found = cache.Get("c")
	assert.True(t, found, "Expected to find key 'c'")
	assert.Equal(t, 3, value, "Expected value 3 for key 'c'")
}

func TestLRUCache_EmptyCache(t *testing.T) {
	cache := NewLRUCache(2)

	// Try getting an item from an empty cache
	_, found := cache.Get("a")
	assert.False(t, found, "Expected not to find 'a' in an empty cache")
}

func TestLRUCache_AddAndEvictMultipleTimes(t *testing.T) {
	cache := NewLRUCache(2)

	// Add and evict multiple times
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3) // Evict "a"

	_, found := cache.Get("a")
	assert.False(t, found, "Expected 'a' to be evicted")

	// Add more items, evicting "b"
	cache.Add("d", 4)

	_, found = cache.Get("b")
	assert.False(t, found, "Expected 'b' to be evicted")

	// Now only "c" and "d" should remain
	value, found := cache.Get("c")
	assert.True(t, found, "Expected to find key 'c'")
	assert.Equal(t, 3, value, "Expected value 3 for key 'c'")

	value, found = cache.Get("d")
	assert.True(t, found, "Expected to find key 'd'")
	assert.Equal(t, 4, value, "Expected value 4 for key 'd'")
}
