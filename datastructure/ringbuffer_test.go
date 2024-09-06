package datastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRingBuffer(t *testing.T) {
	t.Run("size not valid should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewRingBuffer[int](0)
		})
	})

	t.Run("valid size should success", func(t *testing.T) {
		rb := NewRingBuffer[int](10)
		assert.NotNil(t, rb)
		assert.Equal(t, 10, rb.size)
		assert.Equal(t, 0, rb.count)
		assert.Equal(t, 0, rb.head)
		assert.Equal(t, 0, rb.tail)
		assert.Equal(t, 10, len(rb.data))
	})
}

func TestRingBuffer_IsEmpty(t *testing.T) {
	t.Run("empty ring buffer should return true", func(t *testing.T) {
		rb := NewRingBuffer[int](10)
		assert.True(t, rb.IsEmpty())
	})

	t.Run("not empty ring buffer should return false", func(t *testing.T) {
		rb := NewRingBuffer[int](10)
		_ = rb.Enqueue(1)
		assert.False(t, rb.IsEmpty())
	})
}

func TestRingBuffer_IsFull(t *testing.T) {
	t.Run("not full ring buffer should return false", func(t *testing.T) {
		rb := NewRingBuffer[int](10)
		assert.False(t, rb.IsFull())
	})

	t.Run("full ring buffer should return true", func(t *testing.T) {
		rb := NewRingBuffer[int](1)
		_ = rb.Enqueue(1)
		assert.True(t, rb.IsFull())
	})
}

func TestRingBuffer_Enqueue(t *testing.T) {
	t.Run("not override full ring buffer should return error", func(t *testing.T) {
		rb := NewRingBuffer[int](1)
		_ = rb.Enqueue(1)
		err := rb.Enqueue(2)
		assert.Error(t, err)
	})

	t.Run("override full ring buffer should success and cover the oldest", func(t *testing.T) {
		rb := NewRingBuffer[int](2)
		rb.SetOverride(true)
		_ = rb.Enqueue(1)
		_ = rb.Enqueue(2)
		assert.NoError(t, rb.Enqueue(3))
		assert.Equal(t, []int{2, 3}, rb.Data())
	})
}

func TestRingBuffer_Dequeue(t *testing.T) {
	t.Run("empty ring buffer should return not exists", func(t *testing.T) {
		rb := NewRingBuffer[int](1)
		_, exists := rb.Dequeue()
		assert.False(t, exists)
	})
}

func TestRingBuffer_Data(t *testing.T) {
	t.Run("empty ring buffer should return empty slice", func(t *testing.T) {
		rb := NewRingBuffer[int](1)
		assert.Equal(t, []int{}, rb.Data())
	})

	t.Run("not full ring buffer should just return correct data", func(t *testing.T) {
		rb := NewRingBuffer[int](10)
		_ = rb.Enqueue(1)
		_ = rb.Enqueue(2)
		_ = rb.Enqueue(3)
		assert.Equal(t, []int{1, 2, 3}, rb.Data())
		assert.Equal(t, 3, rb.Count())
		assert.Equal(t, 10, rb.Capacity())
	})

	t.Run("circular ring buffer should return data in right order", func(t *testing.T) {
		rb := NewRingBuffer[int](10)
		_ = rb.Enqueue(1)
		_ = rb.Enqueue(2)
		_ = rb.Enqueue(3)
		_ = rb.Enqueue(4)
		_ = rb.Enqueue(5)
		_ = rb.Enqueue(6)
		_ = rb.Enqueue(7)
		_ = rb.Enqueue(8)
		_ = rb.Enqueue(9)
		_ = rb.Enqueue(10)
		v, exists := rb.Dequeue()
		assert.True(t, exists)
		assert.Equal(t, 1, v)
		_ = rb.Enqueue(11)
		assert.Equal(t, []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, rb.Data())
	})
}

func TestRingBuffer_PeekHead(t *testing.T) {
	t.Run("empty ring buffer should return false", func(t *testing.T) {
		rb := NewRingBuffer[int](1)
		_, exists := rb.PeekHead()
		assert.False(t, exists)
	})

	t.Run("not empty ring buffer should return the oldest data", func(t *testing.T) {
		rb := NewRingBuffer[int](3)
		_ = rb.Enqueue(1)
		_ = rb.Enqueue(2)
		_ = rb.Enqueue(3)
		_, _ = rb.Dequeue()
		_ = rb.Enqueue(4)
		v, exists := rb.PeekHead()
		assert.True(t, exists)
		assert.Equal(t, 2, v)
	})
}

func TestRingBuffer_PeekTail(t *testing.T) {
	t.Run("empty ring buffer should return false", func(t *testing.T) {
		rb := NewRingBuffer[int](1)
		_, exists := rb.PeekTail()
		assert.False(t, exists)
	})

	t.Run("not empty ring buffer should return the newest data", func(t *testing.T) {
		rb := NewRingBuffer[int](3)
		_ = rb.Enqueue(1)
		_ = rb.Enqueue(2)
		v, exists := rb.PeekTail()
		assert.True(t, exists)
		assert.Equal(t, 2, v)
	})
}
