package datastructure

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeQueue_Basic(t *testing.T) {
	q := NewLockFreeQueue[int](4)

	// 测试基本的Push操作
	if ok := q.Push(1); !ok {
		t.Error("Push should succeed")
	}
	if ok := q.Push(2); !ok {
		t.Error("Push should succeed")
	}

	// 测试Pop操作
	if val, ok := q.Pop(); !ok || val != 1 {
		t.Errorf("Pop should return 1, got %v", val)
	}
	if val, ok := q.Pop(); !ok || val != 2 {
		t.Errorf("Pop should return 2, got %v", val)
	}

	// 测试空队列的Pop操作
	if _, ok := q.Pop(); ok {
		t.Error("Pop from empty queue should return false")
	}
}

func TestLockFreeQueue_Concurrent(t *testing.T) {
	q := NewLockFreeQueue[int](16)
	const numProducers = 4
	const numConsumers = 4
	const itemsPerProducer = 1000

	var wg sync.WaitGroup
	wg.Add(numProducers + numConsumers)

	// 启动生产者
	for p := 0; p < numProducers; p++ {
		go func() {
			defer wg.Done()
			for i := 0; i < itemsPerProducer; i++ {
				for !q.Push(i) {
					// 队列满时重试
				}
			}
		}()
	}

	// 用于统计消费的数据
	results := make([]int, numConsumers)
	// 启动消费者
	for c := 0; c < numConsumers; c++ {
		go func(consumerID int) {
			defer wg.Done()
			count := 0
			for count < numProducers*itemsPerProducer/numConsumers {
				if _, ok := q.Pop(); ok {
					results[consumerID]++
					count++
				}
			}
		}(c)
	}

	wg.Wait()

	// 验证所有数据都被消费
	total := 0
	for _, count := range results {
		total += count
	}
	expected := numProducers * itemsPerProducer
	if total != expected {
		t.Errorf("Expected %d items to be consumed, got %d", expected, total)
	}
}

func TestLockFreeQueue_Full(t *testing.T) {
	q := NewLockFreeQueue[int](4)

	// 填满队列
	for i := 0; i < 4; i++ {
		if ok := q.Push(i); !ok {
			t.Errorf("Push %d should succeed", i)
		}
	}

	// 尝试向已满的队列Push
	if ok := q.Push(100); ok {
		t.Error("Push to full queue should return false")
	}

	// Pop一个元素后应该能够Push
	if _, ok := q.Pop(); !ok {
		t.Error("Pop should succeed")
	}
	if ok := q.Push(100); !ok {
		t.Error("Push should succeed after Pop")
	}
}

func TestLockFreeQueue_Size(t *testing.T) {
	assert.Panics(t, func() {
		_ = NewLockFreeQueue[int](3)
	})
}

func BenchmarkLockFreeQueue(b *testing.B) {
	q := NewLockFreeQueue[int](1024)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Push(1)
			q.Pop()
		}
	})
}
