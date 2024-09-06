package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSlidingWindowLog(t *testing.T) {
	t.Run("size not valid should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewSlidingWindowLog(0)
		})
	})

	t.Run("valid size should success", func(t *testing.T) {
		sw := NewSlidingWindowLog(10)
		assert.NotNil(t, sw)
		assert.Equal(t, 10, sw.size)
	})
}

func TestSlidingWindowLog_ShouldWork(t *testing.T) {
	size := 10
	interval := time.Millisecond * time.Duration(size)

	sw := NewSlidingWindowLog(size, interval)

	// first 10 tokens should be allowed
	for i := 0; i < size; i++ {
		if i < size-3 {
			time.Sleep(interval / time.Duration(size)) // just sleep 7/10 interval
		}
		assert.True(t, sw.Allow())
	}

	// in current window, no more tokens should be allowed
	assert.False(t, sw.Allow())

	// sleep for 1/2 interval, some older tokens would be removed, new should be allowed
	time.Sleep(interval / 2)
	assert.True(t, sw.Allow())
}
