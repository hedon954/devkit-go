package cacheevict

import (
	"sync"

	"github.com/hedon954/devkit-go/datastructure"
)

type LRUCache struct {
	mu       sync.Mutex
	capacity int

	// hash contains the cached values key and index mapper
	hash map[string]*datastructure.DoublyLinkedNode[cacheItem]

	// link contains the cached values
	link *datastructure.DoublyLinked[cacheItem]
}

type cacheItem struct {
	key   string
	value any
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		hash:     make(map[string]*datastructure.DoublyLinkedNode[cacheItem], capacity),
		link:     datastructure.NewDoublyLinked[cacheItem](),
	}
}

func (lru *LRUCache) Add(k string, v any) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	node, exists := lru.hash[k]
	if exists {
		node.Value.value = v // update the value
		lru.refresh(k, node)
	} else {
		if lru.link.Count() >= lru.capacity {
			lru.evict()
		}
		lru.put(k, v)
	}
}

func (lru *LRUCache) Get(k string) (any, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.hash[k]; ok {
		lru.refresh(k, node)
		return node.Value.value, true
	}
	return nil, false
}

func (lru *LRUCache) put(k string, v any) {
	lru.hash[k] = lru.link.AddToTail(cacheItem{
		key:   k,
		value: v,
	})
}

func (lru *LRUCache) refresh(k string, node *datastructure.DoublyLinkedNode[cacheItem]) {
	lru.hash[k] = lru.link.MoveToTail(node)
}

func (lru *LRUCache) evict() {
	node := lru.link.RemoveFromHead()
	delete(lru.hash, node.Value.key)
}
