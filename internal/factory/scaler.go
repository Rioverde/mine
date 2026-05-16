package factory

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MinWorkers         = 1
	MaxWorkers         = 10
	ScaleUpThreshold   = 15
	ScaleDownThreshold = 2
	CheckInterval      = 200 * time.Millisecond
)

func autoscaler[T any](ctx context.Context, monitorChan <-chan *T, startWorker, stopWorker func()) {
	timer := time.NewTicker(CheckInterval)
	var workers int64 = MinWorkers

	go func() {
		defer timer.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Info("Nothing to monitor, ctx signal received")
				return
			case <-timer.C:
				queue := len(monitorChan)
				current := atomic.LoadInt64(&workers)
				if queue > ScaleUpThreshold && current < MaxWorkers {
					count := atomic.AddInt64(&workers, 1)
					startWorker()
					slog.Info("Scale UP", "workers", count, "queue", queue, "factory", Name(ctx))
				}
				if queue < ScaleDownThreshold && current > MinWorkers {
					count := atomic.AddInt64(&workers, -1)
					stopWorker()
					slog.Info("Scale DOWN", "workers", count, "queue", queue, "factory", Name(ctx))
				}
			}
		}
	}()
}

func scalableWorker[T, V any](
	ctx context.Context,
	input <-chan *T,
	output chan<- *V,
	quit <-chan struct{},
	wg *sync.WaitGroup,
	transform func(*T) *V,
) {
	defer wg.Done()
	slog.Info("Worker started and ready")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Nothing to monitor, ctx signal received")
			return
		case <-quit:
			slog.Info("Worker stopping: scale down signal")
			return
		case v, ok := <-input:
			if !ok {
				slog.Debug("Worker stopping: input channel closed")
				return
			}
			result := transform(v)
			select {
			case <-ctx.Done():
				slog.Info("Nothing to monitor, ctx signal received")
				return
			case output <- result:
			}
		}
	}
}
