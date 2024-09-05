package ratelimiter

import (
	"sync"
	"time"
)

// FixedWindows implements a fixed window rate limiting algorithm.
// It controls the number of allowed requests within a fixed time interval.
type FixedWindows struct {
	mu          sync.Mutex    // Mutex to protect shared state (count, lastTime, nextWinTime) across multiple goroutines
	size        int           // The maximum number of requests allowed in each fixed window interval
	count       int           // The current count of requests in the current window
	interval    time.Duration // The duration of each fixed window interval
	lastTime    time.Time     // The last time a window reset occurred
	nextWinTime time.Time     // The start time of the next window
}

// NewFixedWindows creates a new FixedWindows rate limiter with a specified size and optional interval.
// The size is the maximum number of requests allowed per interval.
// If no interval is provided, it defaults to 1 second.
func NewFixedWindows(size int, interval ...time.Duration) *FixedWindows {
	now := time.Now()

	fw := &FixedWindows{
		size:     size,
		count:    0,
		interval: time.Second,
		lastTime: now,
	}

	if len(interval) > 0 {
		fw.interval = interval[0]
	}

	fw.nextWinTime = now.Add(fw.interval)
	return fw
}

// Allow checks if a single request can be allowed in the current window.
// It returns true if the request is allowed, and false otherwise.
func (fw *FixedWindows) Allow() bool {
	return fw.AllowN(1)
}

// AllowN checks if 'n' requests can be allowed in the current window.
// It returns true if the requests are allowed, and false otherwise.
func (fw *FixedWindows) AllowN(n int) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	if now.After(fw.nextWinTime) {
		timeWindows := now.Sub(fw.lastTime) / fw.interval
		fw.count = 0
		fw.lastTime = fw.lastTime.Add(timeWindows * fw.interval)
		fw.nextWinTime = fw.lastTime.Add(fw.interval)
	}

	if fw.count+n > fw.size {
		return false // Reject the request(s) if it exceeds the limit
	}

	fw.count += n
	return true
}
