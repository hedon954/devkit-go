package ratelimiter

import (
	"sync"
	"time"

	"github.com/hedon954/devkit-go/datastructure"
)

type SlidingWindowLog struct {
	mu       sync.Mutex
	size     int
	interval time.Duration
	logs     *datastructure.RingBuffer[int64]
}

func NewSlidingWindowLog(size int, interval ...time.Duration) *SlidingWindowLog {
	if size <= 0 {
		panic("size must be greater than 0")
	}

	sw := &SlidingWindowLog{
		size:     size,
		interval: time.Second,
		logs:     datastructure.NewRingBuffer[int64](size),
	}

	if len(interval) > 0 {
		sw.interval = interval[0]
	}

	return sw
}

func (sw *SlidingWindowLog) Allow() bool {
	return sw.AllowN(1)
}

//nolint:all
func (sw *SlidingWindowLog) AllowN(n int) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	return true
}
