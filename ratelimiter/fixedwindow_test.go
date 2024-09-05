package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedWindows_ShouldWork(t *testing.T) {
	const size = 10
	const interval = time.Millisecond

	rl := NewFixedWindows(size, interval)

	// out of size, should be rejected
	assert.False(t, rl.AllowN(size+1))

	// first 10 tokens should be allowed
	for i := 0; i < size; i++ {
		assert.True(t, rl.Allow())
	}

	// in current window, no more tokens should be allowed
	for i := 0; i < size; i++ {
		assert.False(t, rl.Allow())
	}

	// sleep for 1 interval to generate new tokens, should be allowed
	time.Sleep(interval)
	for i := 0; i < size/2; i++ {
		assert.True(t, rl.Allow())
	}
	time.Sleep(interval / 2)
	for i := 0; i < size/2; i++ {
		assert.True(t, rl.Allow())
	}
}
