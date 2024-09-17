//nolint:all
package cacheevict

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIFOCache_AddAndGet(t *testing.T) {
	t.Run("invalid capacity should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewFIFOCache(0)
		})
	})

	t.Run("valid capacity should work", func(t *testing.T) {
		cache := NewFIFOCache(2)

		cache.Add("key1", "value1")
		cache.Add("key2", "value2")

		if val, ok := cache.Get("key1"); !ok || val.(string) != "value1" {
			t.Errorf("expected value1, got %v", val)
		}

		if val, ok := cache.Get("key2"); !ok || val.(string) != "value2" {
			t.Errorf("expected value2, got %v", val)
		}
	})
}

func TestFIFOCache_Overwrite(t *testing.T) {
	cache := NewFIFOCache(2)

	cache.Add("key1", "value1")
	cache.Add("key1", "value1_updated")

	if val, ok := cache.Get("key1"); !ok || val.(string) != "value1_updated" {
		t.Errorf("expected value1_updated, got %v", val)
	}
}

func TestFIFOCache_Eviction(t *testing.T) {
	cache := NewFIFOCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3") // This should evict "key1"

	if _, ok := cache.Get("key1"); ok {
		t.Errorf("expected key1 to be evicted")
	}

	if val, ok := cache.Get("key2"); !ok || val.(string) != "value2" {
		t.Errorf("expected value2, got %v", val)
	}

	if val, ok := cache.Get("key3"); !ok || val.(string) != "value3" {
		t.Errorf("expected value3, got %v", val)
	}
}

func TestFIFOCache_MoveToFront(t *testing.T) {
	cache := NewFIFOCache(2)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key1", "value1_updated") // This should move "key1" to the front

	if val, ok := cache.Get("key1"); !ok || val.(string) != "value1_updated" {
		t.Errorf("expected value1_updated, got %v", val)
	}

	cache.Add("key3", "value3") // This should evict "key2"

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("expected key2 to be evicted")
	}

	if val, ok := cache.Get("key3"); !ok || val.(string) != "value3" {
		t.Errorf("expected value3, got %v", val)
	}
}
