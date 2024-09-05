package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_ShouldWork(t *testing.T) {
	const rate = 10
	const capacity = 100

	rl := NewTokenBucket(rate, capacity)
	rl.SetInterval(time.Millisecond)
	assert.Equal(t, capacity, rl.Tokens())
	assert.Equal(t, capacity, rl.Capacity())
	assert.Equal(t, float64(rate), rl.Rate())
	assert.Equal(t, time.Millisecond, rl.Interval())

	// first 100 tokens should be allowed
	for i := 0; i < capacity; i++ {
		assert.True(t, rl.Allow())
	}

	// next 100 tokens should be rejected
	for i := 0; i < capacity; i++ {
		assert.False(t, rl.Allow())
	}

	// sleep for 1 interval to generate new tokens
	time.Sleep(time.Millisecond)
	for i := 0; i < rate; i++ {
		assert.True(t, rl.Allow())
	}

	// the new tokens have been consumed, new request should be rejected
	for i := 0; i < rate; i++ {
		assert.False(t, rl.Allow())
	}
}
