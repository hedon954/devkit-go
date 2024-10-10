package cacheevict

import (
	"container/list"
	"sync"
)

// ARCCache is a cache using ARC (Adaptive Replacement Cache) algorithm.
// ref: https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf
type ARCCache struct {
	mu                 sync.Mutex
	c                  int
	p                  int
	t1, t2, b1, b2     *list.List
	t1m, t2m, b1m, b2m map[string]*list.Element
}

func NewARCCache(size int) *ARCCache {
	return &ARCCache{
		c:   size,
		p:   0,
		t1:  list.New(),
		b1:  list.New(),
		t2:  list.New(),
		b2:  list.New(),
		t1m: make(map[string]*list.Element),
		b1m: make(map[string]*list.Element),
		t2m: make(map[string]*list.Element),
		b2m: make(map[string]*list.Element),
	}
}

func (c *ARCCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t1.Len() + c.t2.Len()
}

func (c *ARCCache) Get(key string) (any, bool) {
	// TODO: implement ARC
	return nil, false
}

func (c *ARCCache) Add(key string, value any) {
	// TODO: implement ARC
}
