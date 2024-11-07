package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/hedon954/devkit-go/datastructure"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	q := datastructure.NewLockFreeQueue[int64](16)

	num := atomic.Int64{}

	producer := func(id int, q *datastructure.LockFreeQueue[int64]) {
		for {
			sleepTime := time.Duration(rand.Intn(5)) * time.Millisecond // nolint:gosec
			time.Sleep(sleepTime)
			num := num.Add(1)
			if !q.Push(num) {
				fmt.Printf("Producer %d: full\n", id)
			} else {
				fmt.Printf("Producer %d: push %d\n", id, num)
			}
		}
	}

	consumer := func(id int, q *datastructure.LockFreeQueue[int64]) {
		for {
			sleepTime := time.Duration(rand.Intn(10)) * time.Millisecond // nolint:gosec
			time.Sleep(sleepTime)
			value, ok := q.Pop()
			if ok {
				fmt.Printf("Consumer %d: pop %d\n", id, value)
			} else {
				fmt.Printf("Consumer %d: empty\n", id)
			}
		}
	}

	for i := 0; i < 10; i++ {
		go producer(i, q)
		go consumer(i, q)
	}

	time.Sleep(time.Second * 5)
}
