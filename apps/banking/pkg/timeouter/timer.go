package timeouter

import "context"

type Handler func() error

func Run(ctx context.Context, fn Handler) error {
	ch := make(chan error, 1)
	go func() {
		ch <- fn()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}
