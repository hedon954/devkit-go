package cacheevict

import (
	"container/list"
	"sync"
)

// LFUCache implements a cache with the Least Frequently Used (LFU) eviction policy.
type LFUCache struct {
	mu       sync.Mutex
	capacity int
	hash     map[string]*list.Element
	freq     map[int]*list.List
	minFreq  int
}

type lfuEntry struct {
	cacheItem
	freq int
}

// NewLFUCache creates a new LFUCache with the given capacity.
// It panics if the capacity is less than or equal to 0.
func NewLFUCache(capacity int) *LFUCache {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	return &LFUCache{
		capacity: capacity,
		hash:     make(map[string]*list.Element, capacity),
		freq:     make(map[int]*list.List),
		minFreq:  0,
	}
}

// Add inserts a key-value pair into the cache.
// If the key already exists, it updates the value and increments the frequency.
// If the cache is full, it evicts the least frequently used item.
func (c *LFUCache) Add(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.hash[key]; ok {
		c.update(elem, value)
		return
	}

	if len(c.hash) >= c.capacity {
		c.evict()
	}

	entry := &lfuEntry{
		cacheItem: cacheItem{
			key:   key,
			value: value,
		},
		freq: 1,
	}

	if c.freq[1] == nil {
		c.freq[1] = list.New()
	}
	elem := c.freq[1].PushFront(entry)
	c.hash[key] = elem
	c.minFreq = 1
}

// Get retrieves the value for a given key from the cache.
// It returns the value and true if the key exists, otherwise it returns nil and false.
func (c *LFUCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.hash[key]; ok {
		entry := elem.Value.(*lfuEntry)
		c.incrementFreq(entry)
		return entry.value, true
	}

	return nil, false
}

// update updates the value of an existing entry and increments its frequency.
func (c *LFUCache) update(elem *list.Element, value any) {
	entry := elem.Value.(*lfuEntry)
	entry.value = value
	c.incrementFreq(entry)
}

// incrementFreq increments the frequency of an entry and moves it to the appropriate frequency list.
func (c *LFUCache) incrementFreq(entry *lfuEntry) {
	// Remove from current frequency list
	freq := entry.freq
	c.freq[freq].Remove(c.hash[entry.key])
	if c.freq[freq].Len() == 0 {
		delete(c.freq, freq)
		if freq == c.minFreq {
			c.minFreq++
		}
	}

	// Add to new frequency list
	entry.freq++
	if c.freq[entry.freq] == nil {
		c.freq[entry.freq] = list.New()
	}
	elem := c.freq[entry.freq].PushFront(entry)
	c.hash[entry.key] = elem
}

// evict removes the least frequently used entry from the cache.
func (c *LFUCache) evict() {
	// Get the last element of the min frequency l
	l := c.freq[c.minFreq]
	elem := l.Back()
	if elem != nil {
		// Delete from hash map
		entry := elem.Value.(*lfuEntry)
		delete(c.hash, entry.key)
		// Delete from frequency list
		l.Remove(elem)
		if l.Len() == 0 {
			delete(c.freq, c.minFreq)
		}
	}
}
