package cacheevict

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFUCache_NewLFUCache(t *testing.T) {
	t.Run("invalid capacity should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewLFUCache(0)
		})
	})

	t.Run("valid capacity should work", func(t *testing.T) {
		cache := NewLFUCache(2)
		if cache.capacity != 2 {
			t.Errorf("expected capacity 2, got %d", cache.capacity)
		}
		if len(cache.hash) != 0 {
			t.Errorf("expected empty hash map, got %d elements", len(cache.hash))
		}
		if len(cache.freq) != 0 {
			t.Errorf("expected empty freq map, got %d elements", len(cache.freq))
		}
	})
}

func TestLFUCache_AddAndGet(t *testing.T) {
	cache := NewLFUCache(2)

	cache.Add("a", 1)
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("expected to get 1, got %v", val)
	}

	cache.Add("b", 2)
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Errorf("expected to get 2, got %v", val)
	}

	cache.Add("c", 3)
	if _, ok := cache.Get("a"); ok {
		t.Errorf("expected 'a' to be evicted")
	}
	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Errorf("expected to get 3, got %v", val)
	}
}

func TestLFUCache_Eviction(t *testing.T) {
	cache := NewLFUCache(2)

	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Add("b", 4)
	cache.Add("c", 3)

	if _, ok := cache.Get("a"); ok {
		t.Errorf("expected 'ac' to be evicted")
	}
	if val, ok := cache.Get("b"); !ok || val != 4 {
		t.Errorf("expected to get 4, got %v", val)
	}
	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Errorf("expected to get c, got %v", val)
	}
}

func TestLFUCache_UpdateFrequency(t *testing.T) {
	cache := NewLFUCache(2)

	cache.Add("a", 1)
	cache.Add("b", 2)
	cache.Get("a")
	cache.Get("a")
	cache.Add("c", 3)

	if _, ok := cache.Get("b"); ok {
		t.Errorf("expected 'b' to be evicted")
	}
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("expected to get 1, got %v", val)
	}
	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Errorf("expected to get 3, got %v", val)
	}
}
