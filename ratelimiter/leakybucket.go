package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	mu           sync.Mutex
	rate         int
	capacity     int
	currentLevel int
	interval     time.Duration
	queue        chan chan struct{}
	notifyPool   sync.Pool
}

func NewLeakyBucket(rate, capacity int, interval ...time.Duration) *LeakyBucket {
	l := &LeakyBucket{
		rate:         rate,
		capacity:     capacity,
		currentLevel: 0,
		interval:     time.Second,
		queue:        make(chan chan struct{}, capacity),
		notifyPool: sync.Pool{
			New: func() interface{} {
				return make(chan struct{})
			},
		},
	}

	if len(interval) > 0 {
		l.interval = interval[0]
	}

	go l.start()

	return l
}

func (l *LeakyBucket) start() {
	ticker := time.NewTicker(l.interval)
	for range ticker.C {
		for i := 0; i < l.rate; i++ {
			notify, ok := <-l.queue
			if !ok {
				break
			}
			notify <- struct{}{}
		}
	}
}

func (l *LeakyBucket) Allow() bool {
	if !l.try() {
		return false
	}

	ok := l.notifyPool.Get().(chan struct{})
	l.queue <- ok

	<-ok
	l.leak()
	l.notifyPool.Put(ok)

	return true
}

func (l *LeakyBucket) try() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.currentLevel >= l.capacity {
		return false
	}
	l.currentLevel++
	return true
}

func (l *LeakyBucket) leak() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.currentLevel > 0 {
		l.currentLevel--
	}
}

func (l *LeakyBucket) Rate() int {
	return l.rate
}

func (l *LeakyBucket) Capacity() int {
	return l.capacity
}

func (l *LeakyBucket) Interval() time.Duration {
	return l.interval
}
