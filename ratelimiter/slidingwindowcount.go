package ratelimiter

import (
	"sync"
	"time"
)

// SlidingWindowCount represents a rate limiter based on the sliding window count algorithm.
type SlidingWindowCount struct {
	mu             sync.Mutex
	buckets        []int         // Number of requests in each time bucket
	size           int           // Maximum allowed requests in the window
	interval       time.Duration // Total sliding window size (e.g., 1 second)
	bucketInterval time.Duration // Size of each time bucket (e.g., 100 milliseconds)
	lastTime       time.Time     // The last time when a request was made
	lastIndex      int           // Index of the last bucket that was updated
}

// NewSlidingWindowCount creates a new Sliding Window Count rate limiter.
// 'size' defines the maximum number of requests allowed within the 'interval' time window.
// 'bucketCount' defines how many buckets the window should be divided into.
func NewSlidingWindowCount(size int, interval time.Duration, bucketCount int) *SlidingWindowCount {
	if size <= 0 || interval <= 0 || bucketCount <= 0 {
		panic("size, interval, and bucketCount must be greater than 0")
	}

	// Calculate the size of each bucket (how long each bucket represents in time)
	bucketSize := interval / time.Duration(bucketCount)

	// Initialize empty buckets
	buckets := make([]int, bucketCount)

	return &SlidingWindowCount{
		buckets:        buckets,
		size:           size,
		interval:       interval,
		bucketInterval: bucketSize,
		lastTime:       time.Now(),
		lastIndex:      0,
	}
}

// Allow checks if a single request is allowed based on the current sliding window.
func (sw *SlidingWindowCount) Allow() bool {
	return sw.AllowN(1)
}

// AllowN checks if 'n' requests are allowed within the current sliding window.
func (sw *SlidingWindowCount) AllowN(n int) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	// Update the buckets based on the current time
	sw.updateBuckets()

	// Check if we can allow the new requests
	if sw.totalCount()+n <= sw.size {
		// Add the new requests to the current time bucket
		sw.addRequests(n)
		return true
	}
	return false
}

// updateBuckets shifts and clears outdated buckets based on the current time.
func (sw *SlidingWindowCount) updateBuckets() {
	now := time.Now()
	bucketPassed := sw.bucketPassed(now)

	for i := 0; i < bucketPassed; i++ {
		idx := (i + sw.lastIndex) % len(sw.buckets)
		sw.buckets[idx] = 0
	}

	sw.lastTime = now
	sw.lastIndex = (sw.lastIndex + bucketPassed) % len(sw.buckets)
}

// totalCount returns the total number of requests in the current sliding window.
func (sw *SlidingWindowCount) totalCount() int {
	total := 0
	for _, count := range sw.buckets {
		total += count
	}
	return total
}

// addRequests adds the given number of requests to the current time bucket.
func (sw *SlidingWindowCount) addRequests(n int) {
	sw.buckets[sw.lastIndex] += n
}

func (sw *SlidingWindowCount) bucketPassed(now time.Time) int {
	res := int(now.Sub(sw.lastTime) / sw.bucketInterval)
	if res > len(sw.buckets) {
		res = len(sw.buckets)
	}
	return res
}
