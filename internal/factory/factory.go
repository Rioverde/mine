package factory

import (
	"context"
	"log/slog"
	"sync"
)

const OutBufferSize = 10

type nameKey struct{}

func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameKey{}, name)
}

func Name(ctx context.Context) string {
	if v, ok := ctx.Value(nameKey{}).(string); ok {
		return v
	}
	return ""
}

func New[T, V any](ctx context.Context, materials <-chan *T, transform func(*T) *V) <-chan *V {
	out := make(chan *V, OutBufferSize)
	quit := make(chan struct{})
	var wg sync.WaitGroup

	start := func() {
		wg.Add(1)
		go scalableWorker(ctx, materials, out, quit, &wg, transform)
	}

	stop := func() {
		select {
		case quit <- struct{}{}:
		default:
		}
	}

	for range MinWorkers {
		start()
	}

	autoscaler(ctx, materials, start, stop)

	go func() {
		wg.Wait()
		close(out)
		slog.Info("All workers finished and channel closed", "factory", Name(ctx))
	}()

	return out
}
