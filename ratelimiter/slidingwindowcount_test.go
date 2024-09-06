package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSlidingWindowCount(t *testing.T) {
	t.Run("size not valid should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewSlidingWindowCount(0, time.Second, 10)
		})
	})

	t.Run("interval not valid should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewSlidingWindowCount(10, 0, 10)
		})
	})

	t.Run("bucket count not valid should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewSlidingWindowCount(10, time.Second, 0)
		})
	})

	t.Run("valid args should success", func(t *testing.T) {
		sw := NewSlidingWindowCount(10, time.Second, 10)
		assert.NotNil(t, sw)
		assert.Equal(t, 10, sw.size)
		assert.Equal(t, time.Second, sw.interval)
		assert.Equal(t, 10, len(sw.buckets))
	})
}

func TestSlidingWindowCount_ShouldWork(t *testing.T) {
	size := 20
	bucketCount := 10
	windowInterval := time.Millisecond * time.Duration(bucketCount)

	sw := NewSlidingWindowCount(size, windowInterval, bucketCount)

	// first 20 tokens should be allowed
	for i := 0; i < size; i++ {
		assert.True(t, sw.Allow())
	}

	// in current window, no more tokens should be allowed
	assert.False(t, sw.Allow())
	assert.Equal(t, size, sw.totalCount())

	// sleep for 1/2 interval, some older tokens would be removed, new should be allowed
	time.Sleep(windowInterval / 2)
	assert.True(t, sw.Allow())

	// sleep for a long time, all buckets should be cleared
	time.Sleep(windowInterval * 2)
	assert.True(t, sw.Allow())
	assert.Equal(t, 1, sw.totalCount())
}
