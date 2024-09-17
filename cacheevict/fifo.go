package cacheevict

import (
	"container/list"
	"sync"
)

// FIFOCache represents a thread-safe FIFO (First-In-First-Out) cache.
type FIFOCache struct {
	mu       sync.RWMutex
	capacity int
	count    int
	hash     map[string]*list.Element
	list     *list.List
}

// NewFIFOCache creates a new FIFOCache with the given capacity.
// It panics if the capacity is less than or equal to 0.
func NewFIFOCache(capacity int) *FIFOCache {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	return &FIFOCache{
		capacity: capacity,
		hash:     make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Add adds a key-value pair to the cache. If the key already exists,
// it updates the value and moves the item to the front of the list.
// If the cache is at capacity, it evicts the oldest item.
func (c *FIFOCache) Add(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if exists, update value and move to front
	if elem, exists := c.hash[key]; exists {
		elem.Value.(*cacheItem).value = value
		c.list.MoveToFront(elem)
		return
	}

	// if out of capacity, remove the last element
	if c.count >= c.capacity {
		backElem := c.list.Back()
		backItem := backElem.Value.(*cacheItem)
		delete(c.hash, backItem.key)
		c.list.Remove(backElem)
		c.count--
	}

	// add the new element to the front
	item := &cacheItem{
		key:   key,
		value: value,
	}
	elem := c.list.PushFront(item)
	c.hash[key] = elem
	c.count++
}

// Get retrieves the value associated with the given key from the cache.
// It returns the value and a boolean indicating whether the key was found.
func (c *FIFOCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if v, ok := c.hash[key]; ok {
		return v.Value.(*cacheItem).value, true
	}
	return nil, false
}
