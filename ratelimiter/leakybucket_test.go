package ratelimiter

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestLeakyBucket_RateLimiting tests the rate limiting functionality
func TestLeakyBucket_RateLimiting(t *testing.T) {
	rate := 5                          // Allows 5 requests per second
	capacity := 10                     // Bucket capacity is 10
	interval := 200 * time.Millisecond // Leaks once every 200ms
	bucket := NewLeakyBucket(rate, capacity, interval)

	var wg sync.WaitGroup
	var processed atomic.Int64

	start := time.Now()

	// Simulate request rate
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if bucket.Allow() {
				processed.Add(1)
			} else {
				t.Errorf("Request rejected")
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	if math.Abs(float64(elapsed.Milliseconds()-10/int64(rate)*interval.Milliseconds())) > 5 {
		t.Errorf("Expected process %dms, got %dms, out of 5ms diff",
			10/int64(rate)*interval.Milliseconds(), elapsed.Milliseconds())
	}
}

// TestLeakyBucket_CapacityLimit tests the capacity limit functionality
func TestLeakyBucket_CapacityLimit(t *testing.T) {
	rate := 1                          // Allows 1 request per second
	capacity := 3                      // Bucket capacity is 3
	interval := 100 * time.Millisecond // Leaks once every 100ms
	bucket := NewLeakyBucket(rate, capacity, interval)

	var wg sync.WaitGroup
	var allowed atomic.Int64

	// Simultaneously make 5 requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if bucket.Allow() {
				allowed.Add(1)
			}
		}()
	}

	wg.Wait()

	// Check if the actual number of allowed requests matches the bucket capacity
	if int(allowed.Load()) != capacity {
		t.Errorf("Allowed more requests than capacity: got %d, want %d", allowed.Load(), capacity)
	}
}

// TestLeakyBucket_RequestRejection tests if requests are correctly rejected when exceeding capacity
func TestLeakyBucket_RequestRejection(t *testing.T) {
	rate := 1                          // Allows 1 request per second
	capacity := 3                      // Bucket capacity is 3
	interval := 100 * time.Millisecond // Leaks once every 100ms
	bucket := NewLeakyBucket(rate, capacity, interval)

	var wg sync.WaitGroup
	var allowed atomic.Int64
	var rejected atomic.Int64

	// Simultaneously make 5 requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if bucket.Allow() {
				allowed.Add(1)
			} else {
				rejected.Add(1)
			}
		}()
	}

	wg.Wait()

	// Check if the actual number of allowed requests matches the bucket capacity
	if int(allowed.Load()) != capacity {
		t.Errorf("Allowed requests did not match capacity: got %d, want %d", allowed.Load(), capacity)
	}

	// Check if the number of rejected requests is correct (should be 2)
	expectedRejected := 5 - capacity
	if int(rejected.Load()) != expectedRejected {
		t.Errorf("Rejected requests did not match expected: got %d, want %d", rejected.Load(), expectedRejected)
	}
}
