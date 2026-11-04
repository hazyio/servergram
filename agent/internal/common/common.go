package common

import "sync"

type Result[T any] struct {
	Val T
	Err error
}

func RunAll[T any](fns ...func() (T, error)) []Result[T] {
	out := make([]Result[T], len(fns))
	var wg sync.WaitGroup
	wg.Add(len(fns))

	for i, fn := range fns {
		go func() {
			defer wg.Done()
			v, err := fn()
			out[i] = Result[T]{v, err}
		}()
	}
	wg.Wait()
	return out
}
