package parallel_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hedon954/devkit-go/parallel"
)

func TestSyncFuncWithReturn(t *testing.T) {
	// Test case 1: Basic functionality
	t.Run("Basic Functionality", func(t *testing.T) {
		ctx := context.Background()
		fn := func() (int, error) {
			return 42, nil
		}

		resultFunc := parallel.RunAsyncWithReturn(ctx, fn)
		result, err := resultFunc()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result != 42 {
			t.Fatalf("expected result 42, got %v", result)
		}
	})

	// Test case 2: Error handling
	t.Run("Error Handling", func(t *testing.T) {
		ctx := context.Background()
		fn := func() (int, error) {
			return 0, errors.New("test error")
		}

		resultFunc := parallel.RunAsyncWithReturn(ctx, fn)
		_, err := resultFunc()

		if err == nil || err.Error() != "test error" {
			t.Fatalf("expected error 'test error', got %v", err)
		}
	})

	// Test case 3: Context cancellation
	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		fn := func() (int, error) {
			time.Sleep(2 * time.Second)
			return 0, nil
		}

		resultFunc := parallel.RunAsyncWithReturn(ctx, fn)
		cancel() // Cancel the context immediately

		_, err := resultFunc()
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled error, got %v", err)
		}
	})

	// Test case 4: Multiple tasks execution time
	t.Run("Multiple Tasks Execution Time", func(t *testing.T) {
		ctx := context.Background()
		fn1 := func() (int, error) {
			time.Sleep(100 * time.Millisecond)
			return 1, nil
		}
		fn2 := func() (int, error) {
			time.Sleep(200 * time.Millisecond)
			return 2, nil
		}

		start := time.Now()
		resultFunc1 := parallel.RunAsyncWithReturn(ctx, fn1)
		resultFunc2 := parallel.RunAsyncWithReturn(ctx, fn2)

		_, _ = resultFunc1()
		_, _ = resultFunc2()
		duration := time.Since(start)

		if duration < 200*time.Millisecond {
			t.Fatalf("expected execution time to be at least 200 milliseconds, got %v", duration)
		}

		if duration >= 300*time.Millisecond {
			t.Fatalf("expected execution time to be between 200 and 300 milliseconds, got %v", duration)
		}
	})

	// Test case 5: Panic handling
	t.Run("Panic Handling", func(t *testing.T) {
		ctx := context.Background()
		fn := func() (int, error) {
			panic("test panic")
		}

		var resultFunc func() (int, error)
		assert.NotPanics(t, func() {
			resultFunc = parallel.RunAsyncWithReturn(ctx, fn)
		})

		assert.Panics(t, func() {
			_, _ = resultFunc()
		})
	})
}
