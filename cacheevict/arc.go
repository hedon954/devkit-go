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
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.lookup(key)
	if !ok {
		return nil, false
	}
	return el.Value.(*cacheItem).value, true
}

func (c *ARCCache) Add(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, _ := c.lookup(key)
	if el != nil {
		el.Value.(*cacheItem).value = value
		return
	}

	// case 4, not in cache and not in ghost, add new item
	l1 := c.t1.Len() + c.b1.Len()
	if l1 >= c.c {
		if c.t1.Len() < c.c {
			el := c.b1.Back()
			c.b1.Remove(el)
			delete(c.b1m, el.Value.(*cacheItem).key)
			c.replacextp(key)
		} else {
			el := c.t1.Back()
			c.t1.Remove(el)
			delete(c.t1m, el.Value.(*cacheItem).key)
		}
	} else {
		if c.t1.Len()+c.t2.Len()+c.b1.Len()+c.b2.Len() >= c.c {
			if c.t1.Len()+c.t2.Len()+c.b1.Len()+c.b2.Len() == 2*c.c {
				el := c.b2.Back()
				c.b2.Remove(el)
				delete(c.b2m, el.Value.(*cacheItem).key)
			}
			c.replacextp(key)
		}
	}
	c.t1.PushFront(&cacheItem{key: key, value: value})
	c.t1m[key] = c.t1.Front()
}

func (c *ARCCache) lookup(key string) (*list.Element, bool) {
	// case 1-1, in t1, move to t2 MRU
	if el, ok := c.t1m[key]; ok {
		c.t1.Remove(el)
		delete(c.t1m, key)
		c.t2.PushFront(el.Value)
		c.t2m[key] = c.t2.Front()
		return el, true
	}

	// case 1-2, in t2, move to t2 MRU
	if el, ok := c.t2m[key]; ok {
		c.t2.MoveToFront(el)
		return el, true
	}

	// case 2, in b1, update p for t1, replace xtp and move to t2 MRU
	if el, ok := c.b1m[key]; ok {
		c.updatePForT1()
		c.replacextp(key)
		c.b1.Remove(el)
		delete(c.b1m, key)
		c.t2.PushFront(el.Value)
		c.t2m[key] = c.t2.Front()
		return el, false
	}

	// case 3, in b2, update p for t2, replace xtp and move to t2 MRU
	if el, ok := c.b2m[key]; ok {
		c.updatePForT2()
		c.replacextp(key)
		c.b2.Remove(el)
		delete(c.b2m, key)
		c.t2.PushFront(el.Value)
		c.t2m[key] = c.t2.Front()
		return el, false
	}

	return nil, false
}

func (c *ARCCache) updatePForT1() {
	c.p = min(c.p+c.δ1(), c.c)
}

func (c *ARCCache) updatePForT2() {
	c.p = max(c.p-c.δ2(), 0)
}

func (c *ARCCache) δ1() int {
	if c.b1.Len() >= c.b2.Len() {
		return 1
	}
	return c.b2.Len() / c.b1.Len()
}

func (c *ARCCache) δ2() int {
	if c.b2.Len() >= c.b1.Len() {
		return 1
	}
	return c.b1.Len() / c.b2.Len()
}

func (c *ARCCache) replacextp(key string) {
	if (c.t1.Len() > 0) && (c.t1.Len() > c.p || (c.b2m[key] != nil && c.t1.Len() == c.p)) {
		// delete LRU from t1 and move to b1 MRU
		el := c.t1.Back()
		c.t1.Remove(el)
		delete(c.t1m, el.Value.(*cacheItem).key)
		c.b1.PushFront(el.Value)
		c.b1m[el.Value.(*cacheItem).key] = c.b1.Front()
	} else {
		// delete LRU from t2 and move to b2 MRU
		el := c.t2.Back()
		c.t2.Remove(el)
		delete(c.t2m, el.Value.(*cacheItem).key)
		c.b2.PushFront(el.Value)
		c.b2m[el.Value.(*cacheItem).key] = c.b2.Front()
	}
}
