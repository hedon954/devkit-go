package ratelimiter

import (
	"sync"
	"time"
)

// LeakyBucket represents a rate limiter using the leaky bucket algorithm.
// It controls the rate at which requests are allowed, ensuring they do not exceed the specified rate and capacity.
type LeakyBucket struct {
	mu           sync.Mutex         // Mutex to protect shared state (currentLevel) across multiple goroutines
	rate         int                // The maximum number of requests allowed per interval
	capacity     int                // The maximum number of requests the bucket can hold at any given time
	currentLevel int                // The current number of requests in the bucket
	interval     time.Duration      // The time interval at which the bucket leaks requests
	queue        chan chan struct{} // A channel of channels to manage request notifications and their order
	notifyPool   sync.Pool          // A pool to reuse channels, reducing the overhead of creating new channels frequently
}

// NewLeakyBucket creates a new LeakyBucket instance with a specified rate, capacity, and optional interval.
// The rate is the number of allowed requests per interval, and the capacity is the maximum number of requests the bucket can hold.
// If no interval is provided, it defaults to 1 second.
func NewLeakyBucket(rate, capacity int, interval ...time.Duration) *LeakyBucket {
	l := &LeakyBucket{
		rate:         rate,
		capacity:     capacity,
		currentLevel: 0,
		interval:     time.Second,                        // Default interval to 1 second if not specified
		queue:        make(chan chan struct{}, capacity), // Buffered channel to handle up to 'capacity' requests
		notifyPool: sync.Pool{
			New: func() interface{} {
				return make(chan struct{}) // Create a new channel when needed
			},
		},
	}

	// Override the default interval if provided
	if len(interval) > 0 {
		l.interval = interval[0]
	}

	// Start the goroutine that will leak requests at a fixed rate
	go l.start()

	return l
}

// start begins a goroutine that regularly leaks requests from the bucket at the defined interval and rate.
func (l *LeakyBucket) start() {
	ticker := time.NewTicker(l.interval)
	for range ticker.C {
		for i := 0; i < l.rate; i++ {
			select {
			case notify := <-l.queue:
				notify <- struct{}{}
			default:
			}
		}
	}
}

// Allow checks if a new request is allowed under the current rate and capacity constraints.
// It returns true if the request is allowed, and false otherwise.
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

// try attempts to increment the bucket's current level and checks if the request can be accommodated.
// It returns true if the request is accepted, false if the bucket is full.
func (l *LeakyBucket) try() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.currentLevel >= l.capacity {
		return false
	}
	l.currentLevel++
	return true
}

// leak decreases the current number of requests in the bucket by one.
// This is called after a request has been processed and space in the bucket is freed.
func (l *LeakyBucket) leak() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.currentLevel > 0 {
		l.currentLevel--
	}
}
