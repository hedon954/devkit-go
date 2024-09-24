package cacheevict

import (
	"container/list"
	"sync"
)

// ARCCache is a cache using ARC (Adaptive Replacement Cache) algorithm.
// ref: https://www.usenix.org/legacy/events/fast03/tech/full_papers/megiddo/megiddo.pdf
type ARCCache struct {
	mu  sync.Mutex
	c   int
	p   int
	t1  *list.List
	b1  *list.List
	t2  *list.List
	b2  *list.List
	t1m map[string]*list.Element
	b1m map[string]*list.Element
	t2m map[string]*list.Element
	b2m map[string]*list.Element
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

func (c *ARCCache) Add(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := c.lookup(key)
	if item != nil {
		item.Value.(*cacheItem).value = value
		return
	}

	c.t1.PushFront(&cacheItem{key: key, value: value})
	c.t1m[key] = c.t1.Front()
}

func (c *ARCCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := c.lookup(key)
	if item == nil {
		return nil, false
	}
	return item.Value.(*cacheItem).value, true
}

func (c *ARCCache) lookup(key string) *list.Element {
	// 1. if the item is in t1 or t2, it's moved to the MRU position of T2
	if item := c.inT1(key); item != nil {
		c.fromT1ToT2(key, item)
		return item
	}
	if item := c.inT2(key); item != nil {
		c.t2.MoveToFront(item)
		return item
	}
	// 2. if the item is in B1, it indicates a recency preference,
	//    and p is adjusted to increase T1's size.
	if item := c.inB1(key); item != nil {
		c.adjustPForB1()
		c.fromB1ToT2(key, item)
		return item
	}
	// 3. if the item is in B2, it indicates a frequency preference,
	//    and p is adjusted to increase T2's size.
	if item := c.inB2(key); item != nil {
		c.adjustPForB2()
		c.fromB2ToT2(key, item)
		return item
	}
	// 4. if the item is not in the cache, it's added to T1.
	// 4.1 t1 + b2 = c
	if c.t1.Len()+c.b2.Len() == c.c {
		//  4.1.1 t1 < c, delete lru page in b1, replace (xt, p)
		if c.t1.Len() < c.c {
			c.evictB1()
			c.replacextp(key)
		} else {
			//  4.1.2 t1 = c, delete lru page in b1
			c.replacextp(key)
		}
	} else if c.t1.Len()+c.b2.Len() < c.c {
		// 4.2 t1 + b2 < c
		// 4.2.1 t1+t2+b1+b2 >= c, delete lru page in b2,
		if c.t1.Len()+c.b1.Len()+c.t2.Len()+c.b2.Len() >= c.c {
			c.evictB2()
			// and if t1+t2+b1+b2 = c, replace (xt, p)
			if c.t1.Len()+c.b1.Len()+c.t2.Len()+c.b2.Len() == c.c {
				c.replacextp(key)
			}
		}
	}
	return nil
}

func (c *ARCCache) inT1(key string) *list.Element {
	return c.t1m[key]
}

func (c *ARCCache) inT2(key string) *list.Element {
	return c.t2m[key]
}

func (c *ARCCache) fromT1ToT2(key string, item *list.Element) {
	c.t1.Remove(item)
	delete(c.t1m, key)
	c.t2.PushFront(item.Value)
	c.t2m[key] = item
}

func (c *ARCCache) fromT1ToB1() {
	item := c.t1.Back()
	c.t1.Remove(item)
	delete(c.t1m, item.Value.(*cacheItem).key)
	c.b1.PushFront(item.Value)
	c.b1m[item.Value.(*cacheItem).key] = item
}

func (c *ARCCache) fromB1ToT2(key string, item *list.Element) {
	c.b1.Remove(item)
	delete(c.b1m, key)
	c.t2.PushFront(item.Value)
	c.t2m[key] = item
}

func (c *ARCCache) fromB2ToT2(key string, item *list.Element) {
	c.b2.Remove(item)
	delete(c.b2m, key)
	c.t2.PushFront(item.Value)
	c.t2m[key] = item
}

func (c *ARCCache) fromT2ToB2() {
	item := c.t2.Back()
	c.t2.Remove(item)
	delete(c.t2m, item.Value.(*cacheItem).key)
	c.b2.PushFront(item.Value)
	c.b2m[item.Value.(*cacheItem).key] = item
}

func (c *ARCCache) evictB1() {
	item := c.b1.Back()
	if item == nil {
		return
	}
	c.b1.Remove(item)
	delete(c.b1m, item.Value.(*cacheItem).key)
}

func (c *ARCCache) evictB2() {
	item := c.b2.Back()
	if item == nil {
		return
	}
	c.b2.Remove(item)
	delete(c.b2m, item.Value.(*cacheItem).key)
}

func (c *ARCCache) inB1(key string) *list.Element {
	return c.b1m[key]
}

func (c *ARCCache) inB2(key string) *list.Element {
	return c.b2m[key]
}

func (c *ARCCache) adjustPForB1() {
	c.p = min(c.c, c.p+c.δ1())
}

func (c *ARCCache) adjustPForB2() {
	c.p = max(c.p-c.δ2(), 0)
}

func (c *ARCCache) δ1() int {
	if c.b1.Len() >= c.b2.Len() {
		return 1
	}
	return int(c.b2.Len() / c.b1.Len())
}

func (c *ARCCache) δ2() int {
	if c.b2.Len() >= c.b1.Len() {
		return 1
	}
	return int(c.b1.Len() / c.b2.Len())
}

func (c *ARCCache) replacextp(key string) {
	if (c.t1.Len() > 0) && ((c.t1.Len() > c.p) || (c.inB2(key) != nil && c.t1.Len() == c.p)) {
		c.fromT1ToB1()
	} else {
		c.fromT2ToB2()
	}
}
