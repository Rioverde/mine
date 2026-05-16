package merchant

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/repo"
)

func Start(ctx context.Context, cfg *config.Config, r repo.Repo) <-chan *repo.ItemDTO {
	out := make(chan *repo.ItemDTO, cfg.Merchant.BufferSize)
	ticker := time.NewTicker(cfg.Merchant.Delay)

	go func() {
		defer ticker.Stop()
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				items, err := r.FetchNextBatch(ctx, cfg.Merchant.BufferSize)
				if err != nil {
					slog.Error("Stream: failed to fetch batch from DB", "err", err)
					continue
				}

				for _, item := range items {
					select {
					case <-ctx.Done():
						return
					case out <- item:
					}
				}
			}
		}
	}()

	return out
}
