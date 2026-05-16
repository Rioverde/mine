package factory

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

func StartMerchantStream(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) <-chan *storage.ItemDTO {
	out := make(chan *storage.ItemDTO, cfg.Merchant.BufferSize)
	ticker := time.NewTicker(cfg.Merchant.Delay)

	go func() {
		defer ticker.Stop()
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				items, err := storage.FetchNextBatch(ctx, pool, cfg.Merchant.BufferSize)
				if err != nil {
					slog.Error("Stream: failed to fetch batch from DB", "err", err)
					return
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
