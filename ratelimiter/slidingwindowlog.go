package ratelimiter

import (
	"slices"
	"sync"
	"time"
)

// SlidingWindowLog is a rate limiter based on the sliding window log algorithm.
// It keeps track of timestamps of requests and allows or rejects new requests
// based on the allowed size and the time window.
type SlidingWindowLog struct {
	mu       sync.Mutex    // Protects concurrent access to logs
	size     int           // Maximum number of allowed requests within the window
	interval time.Duration // The sliding time window (e.g., 1 second)
	logs     []time.Time   // Timestamps of the requests
}

// NewSlidingWindowLog creates a new SlidingWindowLog rate limiter.
// The 'size' defines the maximum number of requests allowed in the 'interval' time window.
// Optionally, a custom 'interval' can be provided. Defaults to 1 second if not specified.
func NewSlidingWindowLog(size int, interval ...time.Duration) *SlidingWindowLog {
	if size <= 0 {
		panic("size must be greater than 0")
	}

	sw := &SlidingWindowLog{
		size:     size,
		interval: time.Second, // Default interval is 1 second
		logs:     make([]time.Time, 0, size),
	}

	// Override interval if provided
	if len(interval) > 0 {
		sw.interval = interval[0]
	}

	return sw
}

// Allow checks if a single request is allowed based on the current window.
func (sw *SlidingWindowLog) Allow() bool {
	return sw.AllowN(1)
}

// AllowN checks if 'n' requests can be allowed within the current time window.
// It first tries to accept the requests if there's enough space without removing old requests.
// If the log is full, it removes expired logs and tries again.
func (sw *SlidingWindowLog) AllowN(n int) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	// If the log is not full, try to accept the new requests directly
	if sw.tryAccept(n, now) {
		return true
	}

	// If the log is full, remove old entries that are outside the window
	sw.removeOlderThan(now.Add(-sw.interval))

	// Try to accept the requests again after removing old entries
	return sw.tryAccept(n, now)
}

// tryAccept attempts to accept 'n' new requests. If there's enough space in the logs,
// the requests will be accepted and timestamps added. Returns true if successful.
func (sw *SlidingWindowLog) tryAccept(n int, now time.Time) bool {
	if len(sw.logs)+n <= sw.size {
		sw.append(n, now)
		return true
	}
	return false
}

func (sw *SlidingWindowLog) append(n int, now time.Time) {
	batch := make([]time.Time, n)
	for i := 0; i < n; i++ {
		batch[i] = now
	}
	sw.logs = append(sw.logs, batch...)
}

func (sw *SlidingWindowLog) removeOlderThan(threshold time.Time) {
	sw.logs = slices.DeleteFunc(sw.logs, func(t time.Time) bool {
		return t.Before(threshold)
	})
}
