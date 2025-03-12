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
	value, found := cache.Get("a")
	if !found || value != 10 {
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
	cache := NewARCCache(4)

	logState := func(msg string) {
		t.Logf("%s: p=%d, t1=%d, t2=%d, b1=%d, b2=%d",
			msg, cache.p, cache.t1.Len(), cache.t2.Len(), cache.b1.Len(), cache.b2.Len())
	}

	// Initial state
	logState("Initial state")

	// Fill the cache
	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("c", 3)
	cache.Add("d", 4)
	logState("After filling cache")

	// Access "a" and "b" to move them to T2
	cache.Get("a")
	cache.Get("b")
	logState("After accessing 'a' and 'b'")

	// Add a new item, causing eviction from T1 to B1
	cache.Add("e", 5)
	logState("After adding 'e'")

	// Access an item in B1, should increase p
	oldP := cache.p
	cache.Get("c")
	logState("After accessing 'c' (B1 hit)")
	if cache.p <= oldP {
		t.Errorf("Expected p to increase after B1 hit, old p: %d, new p: %d", oldP, cache.p)
	}

	// Continue adding new items
	cache.Add("f", 6)
	cache.Add("g", 7)
	logState("After adding 'f' and 'g'")

	// Access an item in B2, should decrease p
	oldP = cache.p
	cache.Get("a")
	logState("After accessing 'a' (B2 hit)")
	if cache.p >= oldP {
		t.Errorf("Expected p to decrease after B2 hit, old p: %d, new p: %d", oldP, cache.p)
	}
}
