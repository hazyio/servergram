package common

import "context"

func RunWithContext[T any](ctx context.Context, fn func() (T, error), cleanup func(T)) (T, error) {
	type result struct {
		val T
		err error
	}
	ch := make(chan result, 1)
	go func() {
		v, err := fn()
		ch <- result{v, err}
	}()

	select {
	case r := <-ch:
		return r.val, r.err
	case <-ctx.Done():
		go func() {
			if r := <-ch; r.err == nil && cleanup != nil {
				cleanup(r.val) // we gave up, so dispose of the late result
			}
		}()
		var zero T
		return zero, ctx.Err()
	}
}
