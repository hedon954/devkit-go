package ratelimiter

import (
	"sync"
	"time"
)

// TokenBucket is a token bucket limiter.
type TokenBucket struct {
	mu sync.Mutex

	// rate is the number of tokens that can be consumed in an interval.
	rate float64

	// capacity is the maximum number of tokens the limiter can hold.
	//
	// In go official implementation, it is called `burst`.
	// `capacity` and `burst` actually are the same thing.
	//
	// The naming difference between burst and capacity is
	// mainly due to different backgrounds and usage scenarios:
	// - burst: Commonly used when describing request processing capacity,
	//			emphasizing the number of requests that can be processed
	//			in a short period of time. For example, in network flow limiting,
	//			burst is used to describe the allowable instantaneous request peak.
	// - capacity: More from the perspective of technical implementation,
	//			it describes the capacity limit of a bucket and represents
	//			the maximum number of stored tokens of a bucket.
	capacity int

	// tokens is the number of available tokens in the bucket.
	tokens float64

	// interval is the time interval of the bucket to generate new tokens.
	interval time.Duration

	// lastTime is the time of the last token generation.
	lastTime time.Time
}

func (l *TokenBucket) SetInterval(i time.Duration) {
	l.interval = i
}

// NewTokenBucket creates a new token bucket limiter.
func NewTokenBucket(rate float64, capacity int) *TokenBucket {
	l := &TokenBucket{
		rate:     rate,
		capacity: capacity,
		tokens:   float64(capacity), // initialize with full capacity
		interval: time.Second,       // default interval is 1 second
		lastTime: time.Now(),        // initialize lastTime to current time
	}

	return l
}

// Allow returns true if the limiter allows 1 token to be processed.
func (l *TokenBucket) Allow() bool {
	return l.AllowN(1)
}

// AllowN returns true if the limiter allows n tokens to be processed.
func (l *TokenBucket) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.advance(time.Now())

	if float64(n) > l.tokens {
		return false
	}

	l.tokens -= float64(n)
	return true
}

// advance advances the limiter to the next time interval,
// and generates new tokens.
func (l *TokenBucket) advance(n time.Time) {
	if l.lastTime.IsZero() {
		l.lastTime = n
		return
	}

	diff := n.Sub(l.lastTime)
	if diff < l.interval {
		return
	}

	// Calculate how many full intervals have passed
	intervalCount := diff / l.interval

	// Add tokens bases on the number of intervals that have passed
	l.tokens += float64(intervalCount) * l.rate
	if l.tokens > float64(l.capacity) {
		l.tokens = float64(l.capacity)
	}

	// Update lastTime to reflect the intervals that have passed
	l.lastTime = l.lastTime.Add(l.interval * intervalCount)
}

func (l *TokenBucket) Tokens() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return int(l.tokens)
}

func (l *TokenBucket) Capacity() int {
	return l.capacity
}

func (l *TokenBucket) Rate() float64 {
	return l.rate
}

func (l *TokenBucket) Interval() time.Duration {
	return l.interval
}

func (l *TokenBucket) LastTime() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.lastTime
}
