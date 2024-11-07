package datastructure

import (
	"sync/atomic"
)

// LockFreeQueue is a lock-free queue.
type LockFreeQueue[T any] struct {
	// head is the next position to push
	head uint64
	// tail is the next position to pop
	tail uint64
	// mod is the mask of the size
	mod uint64
	// size is the size of the queue, must be a power of 2,
	// in order to use the mask to calculate the position,
	// it is faster than the % operation
	size uint64
	// data stores the data of the queue
	data []T
	// valid marks the data is available or not
	valid []bool
}

// NewLockFreeQueue creates a new lock-free queue with the given size.
// The size must be a power of 2 and greater than 0.
// If your producers and consumers are both many, we do not recommend you set the size too small,
// on the one hand, the spinning will be too much and performance will decrease,
// on the other hand, it will increase the possibility of the ABA problem of CAS.
func NewLockFreeQueue[T any](size int) *LockFreeQueue[T] {
	if size <= 0 {
		panic("size must be greater than 0")
	}
	if size&(size-1) != 0 {
		panic("size must be a power of 2")
	}

	return &LockFreeQueue[T]{
		size:  uint64(size),
		mod:   uint64(size) - 1,
		data:  make([]T, size),
		valid: make([]bool, size),
	}
}

// Push is a function that producer use to push data into the queue.
// The thing that producer does is try to grab the head of the queue.
// If grabbed, it will move the head to the next position,
// then store the value into the origin head.
func (q *LockFreeQueue[T]) Push(value T) bool {
	for {
		head := q.head
		if q.valid[head] {
			return false
		}
		// NOTE: here we do not consider the ABA problem,
		// because the head is always moving forward,
		// so the ABA problem is very unlikely to happen.
		// Only when head moves n rounds, it is possible,
		// and it is almost impossible to happen when size is not too small.
		if !atomic.CompareAndSwapUint64(&q.head, head, (head+1)&q.mod) {
			continue
		}
		q.data[head] = value
		q.valid[head] = true
		return true
	}
}

// Pop is a function that consumer use to pop data from the queue.
// The thing that consumer does is try to grab the tail of the queue.
// If grabbed, it will move the tail to the next position,
// then return the value at origin tail and set the valid to false.
func (q *LockFreeQueue[T]) Pop() (T, bool) {
	for {
		tail := q.tail
		if !q.valid[tail] {
			var zero T
			return zero, false
		}
		// NOTE: here we do not consider the ABA problem,
		// the reason is the same as the Push function.
		if !atomic.CompareAndSwapUint64(&q.tail, tail, (tail+1)&q.mod) {
			continue
		}
		q.valid[tail] = false
		return q.data[tail], true
	}
}
