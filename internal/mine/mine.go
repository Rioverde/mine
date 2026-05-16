package mine

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/domain"
)

func Run(ctx context.Context, cfg *config.Config) <-chan *domain.Ore {
	out := make(chan *domain.Ore, cfg.Mine.BufferSize)
	ticker := time.NewTicker(cfg.Mine.Delay)

	go func() {
		defer ticker.Stop()
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				slog.Info("Mine closed")
				return
			case <-ticker.C:
			}

			o := domain.RandomOre(cfg.Ore)
			select {
			case out <- o:
				slog.Info("Ore found", "material", o.Name(), "capacity", o.Capacity)
			case <-ctx.Done():
				slog.Info("Mine closed")
				return
			}
		}
	}()

	return out
}
