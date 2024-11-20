package parallel

import (
	"context"
	"fmt"
	"runtime/debug"
)

// RunAsyncWithReturn runs the given function asynchronously
// and returns a function that can be used to get the result.
func RunAsyncWithReturn[T any](ctx context.Context, fn func() (T, error)) func() (T, error) {
	channel := make(chan T)
	errorChannel := make(chan error)
	panicChannel := make(chan string)

	go func() {
		defer close(channel)
		defer close(errorChannel)
		defer close(panicChannel)
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				select {
				case panicChannel <- fmt.Sprintf("%v\n%s", r, stackTrace):
				case <-ctx.Done():
				}
			}
		}()
		result, err := fn()
		if err != nil {
			select {
			case errorChannel <- err:
			case <-ctx.Done():
			}
		} else {
			select {
			case channel <- result:
			case <-ctx.Done():
			}
		}
	}()

	return func() (T, error) {
		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		case result := <-channel:
			return result, nil
		case err := <-errorChannel:
			var zero T
			return zero, err
		case p := <-panicChannel:
			panic(p)
		}
	}
}
