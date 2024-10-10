package cacheevict

import (
	"testing"
)

func TestARCCache_AddAndGet(t *testing.T) {
	cache := NewARCCache(2)
	cache.Add("a", 1)
	cache.Add("b", 2)

	// Test retrieving existing items
	if value, found := cache.Get("a"); !found || value != 1 {
		t.Errorf("Expected to find key 'a' with value 1, got %v", value)
	}
	if value, found := cache.Get("b"); !found || value != 2 {
		t.Errorf("Expected to find key 'b' with value 2, got %v", value)
	}

	// Test retrieving non-existing item
	if _, found := cache.Get("d"); found {
		t.Errorf("Expected not to find key 'd'")
	}
}

func TestARCCache_UpdateValue(t *testing.T) {
	cache := NewARCCache(3)
	cache.Add("a", 1)

	// Update the value of an existing item
	cache.Add("a", 10)

	// Test that the value is updated
	if value, found := cache.Get("a"); !found || value != 10 {
		t.Errorf("Expected to find key 'a' with updated value 10, got %v", value)
	}
}

func TestARCCache_Eviction(t *testing.T) {
	cache := NewARCCache(3)
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)

	// Add more items to trigger eviction
	cache.Add("d", 4)
	cache.Add("e", 5)

	// Check the state of the cache
	for _, key := range []string{"a", "b", "c", "d", "e"} {
		value, found := cache.Get(key)
		t.Logf("Key %s: found = %v, value = %v", key, found, value)
	}

	// Ensure the cache size is still 3
	if cache.Len() != 3 {
		t.Errorf("Expected cache size to be 3, got %d", cache.Len())
	}

	// Access 'a' again to potentially bring it back from B1 to T2
	cache.Get("a")

	// Add another item
	cache.Add("f", 6)

	// Check the state again
	for _, key := range []string{"a", "b", "c", "d", "e", "f"} {
		value, found := cache.Get(key)
		t.Logf("Key %s: found = %v, value = %v", key, found, value)
	}

	// Ensure the cache size is still 3
	if cache.Len() != 3 {
		t.Errorf("Expected cache size to be 3, got %d", cache.Len())
	}
}

func TestARCCache_AdjustP(t *testing.T) {
	cache := NewARCCache(3)
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)

	// Record initial p value
	initialP := cache.p

	// Add a new item to trigger eviction to B1
	cache.Add("d", 4)

	// Access "a" which should now be in B1, this should increase p
	cache.Get("a")

	if cache.p <= initialP {
		t.Errorf("Expected p to increase after accessing item in B1, got p = %d", cache.p)
	}

	// Add more items to push some to B2
	cache.Add("e", 5)
	cache.Add("f", 6)

	// Record p value after increase
	increasedP := cache.p

	// Access an item that should be in B2, this should decrease p
	cache.Get("c")

	if cache.p >= increasedP {
		t.Errorf("Expected p to decrease after accessing item in B2, got p = %d", cache.p)
	}

	// Final boundary check
	if cache.p < 0 || cache.p > cache.c {
		t.Errorf("p value out of bounds: p = %d, should be 0 <= p <= %d", cache.p, cache.c)
	}
}
