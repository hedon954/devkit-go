package cacheevict

import (
	"testing"
)

func newTestARCCache(size int) *ARCCache {
	cache := NewARCCache(size)
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)
	return cache
}

func TestARCCache_AddAndGet(t *testing.T) {
	cache := newTestARCCache(3)

	// Test retrieving existing items
	if value, found := cache.Get("a"); !found || value != 1 {
		t.Errorf("Expected to find key 'a' with value 1, got %v", value)
	}
	if value, found := cache.Get("b"); !found || value != 2 {
		t.Errorf("Expected to find key 'b' with value 2, got %v", value)
	}
	if value, found := cache.Get("c"); !found || value != 3 {
		t.Errorf("Expected to find key 'c' with value 3, got %v", value)
	}

	// Test retrieving non-existing item
	if _, found := cache.Get("d"); found {
		t.Errorf("Expected not to find key 'd'")
	}
}

func TestARCCache_UpdateValue(t *testing.T) {
	cache := newTestARCCache(3)

	// Update the value of an existing item
	cache.Add("a", 10)

	// Test that the value is updated
	if value, found := cache.Get("a"); !found || value != 10 {
		t.Errorf("Expected to find key 'a' with updated value 10, got %v", value)
	}
}

func TestARCCache_Eviction(t *testing.T) {
	cache := newTestARCCache(3)

	// Add more items to trigger eviction
	cache.Add("d", 4)
	cache.Add("e", 5)

	// Test that the least recently used items are evicted
	if _, found := cache.Get("a"); found {
		t.Errorf("Expected key 'a' to be evicted")
	}
	if _, found := cache.Get("b"); found {
		t.Errorf("Expected key 'b' to be evicted")
	}

	// Test that the most recently used items are still in the cache
	if value, found := cache.Get("c"); !found || value != 3 {
		t.Errorf("Expected to find key 'c' with value 3, got %v", value)
	}
	if value, found := cache.Get("d"); !found || value != 4 {
		t.Errorf("Expected to find key 'd' with value 4, got %v", value)
	}
	if value, found := cache.Get("e"); !found || value != 5 {
		t.Errorf("Expected to find key 'e' with value 5, got %v", value)
	}
}

func TestARCCache_AdjustP(t *testing.T) {
	cache := newTestARCCache(3)

	// Add items to fill the cache and trigger eviction
	cache.Add("d", 4)
	cache.Add("e", 5)

	// Access items in B1 and B2 to adjust p
	cache.Get("a") // Should adjust p to increase T1's size, but return false
	cache.Get("b") // Should adjust p to increase T2's size, but return false

	// Test that p is adjusted correctly
	if cache.p <= 0 || cache.p >= cache.c {
		t.Errorf("Expected p to be adjusted within bounds, got %d", cache.p)
	}
}

func TestARCCache_ARCBenefits(t *testing.T) {
	cache := newTestARCCache(3)

	// Access pattern to test ARC benefits
	cache.Get("a")
	cache.Get("b")
	cache.Get("c")
	cache.Add("d", 4)
	cache.Get("a")
	cache.Add("e", 5)
	cache.Get("b")
	cache.Get("c")
	cache.Get("d")
	cache.Get("e")

	// Test that frequently accessed items are still in the cache
	if value, found := cache.Get("a"); !found || value != 1 {
		t.Errorf("Expected to find key 'a' with value 1, got %v", value)
	}
	if value, found := cache.Get("b"); !found || value != 2 {
		t.Errorf("Expected to find key 'b' with value 2, got %v", value)
	}
	if value, found := cache.Get("c"); !found || value != 3 {
		t.Errorf("Expected to find key 'c' with value 3, got %v", value)
	}
	if value, found := cache.Get("d"); !found || value != 4 {
		t.Errorf("Expected to find key 'd' with value 4, got %v", value)
	}
	if value, found := cache.Get("e"); !found || value != 5 {
		t.Errorf("Expected to find key 'e' with value 5, got %v", value)
	}
}
