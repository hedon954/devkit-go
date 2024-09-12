package datastructure

import (
	"errors"
)

// RingBuffer represents the circular buffer.
type RingBuffer[T any] struct {
	data        []T
	head, tail  int
	size, count int
	override    bool
}

// NewRingBuffer creates a new RingBuffer with a specified size.
func NewRingBuffer[T any](size int) *RingBuffer[T] {
	if size <= 0 {
		panic("size must be greater than 0")
	}
	return &RingBuffer[T]{
		data:     make([]T, size),
		size:     size,
		head:     0,
		tail:     0,
		count:    0,
		override: false,
	}
}

func (rb *RingBuffer[T]) SetOverride(override bool) {
	rb.override = override
}

// IsFull checks if the ring buffer is full.
func (rb *RingBuffer[T]) IsFull() bool {
	return rb.count == rb.size
}

// IsEmpty checks if the ring buffer is empty.
func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.count == 0
}

// Count returns the number of elements in the ring buffer.
func (rb *RingBuffer[T]) Count() int {
	return rb.count
}

// Capacity returns the capacity of the ring buffer.
func (rb *RingBuffer[T]) Capacity() int {
	return rb.size
}

// Enqueue adds an element to the ring buffer.
// Returns an error if the buffer is full and override is disabled.
func (rb *RingBuffer[T]) Enqueue(value T) error {
	if rb.IsFull() {
		// Buffer is full and override is not allowed
		if !rb.override {
			return errors.New("buffer is full")
		}
		// Buffer is full but override is allowed, overwrite the oldest element
		rb.head = (rb.head + 1) % rb.size
	} else {
		// Buffer is not full, increment the count
		rb.count++
	}

	// Insert the new Value
	rb.data[rb.tail] = value
	rb.tail = (rb.tail + 1) % rb.size
	return nil
}

// Dequeue removes and returns the oldest element from the ring buffer.
// Returns an error if the buffer is empty.
func (rb *RingBuffer[T]) Dequeue() (T, bool) {
	var value T
	if rb.IsEmpty() {
		return value, false
	}
	value = rb.data[rb.head]
	rb.head = (rb.head + 1) % rb.size
	rb.count--
	return value, true
}

// PeekHead returns the oldest element in the ring buffer without removing it.
func (rb *RingBuffer[T]) PeekHead() (T, bool) {
	if rb.IsEmpty() {
		var none T
		return none, false
	}
	return rb.data[rb.head], true
}

// PeekTail returns the newest element in the ring buffer without removing it.
func (rb *RingBuffer[T]) PeekTail() (T, bool) {
	if rb.IsEmpty() {
		var none T
		return none, false
	}
	return rb.data[(rb.tail-1+rb.size)%rb.size], true
}

// PrintBuffer prints the current state of the ring buffer with valid data.
func (rb *RingBuffer[T]) Data() []T {
	res := make([]T, rb.count)
	for i := 0; i < rb.count; i++ {
		idx := (rb.head + i) % rb.size
		res[i] = rb.data[idx]
	}
	return res
}
