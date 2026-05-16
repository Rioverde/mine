package factory

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Rioverde/mine/internal/config"
)

func autoscaler[T any](ctx context.Context, cfg config.ScalerConfig, monitorChan <-chan *T, startWorker, stopWorker func()) {
	timer := time.NewTicker(cfg.CheckInterval)
	workers := int64(cfg.MinWorkers)

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
				if queue > cfg.ScaleUpThreshold && current < int64(cfg.MaxWorkers) {
					count := atomic.AddInt64(&workers, 1)
					startWorker()
					slog.Info("Scale UP", "workers", count, "queue", queue, "factory", Name(ctx))
				}
				if queue < cfg.ScaleDownThreshold && current > int64(cfg.MinWorkers) {
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
