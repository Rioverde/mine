package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/repo"
)

func OutboxPublisher(ctx context.Context, r repo.Repo, bus Bus, cfg *config.AppConfig) {
	ticker := time.NewTicker(cfg.CacheSyncDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sent, err := r.DrainOutbox(ctx, cfg.BufferSize, func(row repo.OutboxRow) error {
				_, err := bus.Publish(ctx, row.Payload)
				if err != nil {
					slog.Error("outbox: publish", "err", err, "outbox_id", row.ID, "item_id", row.ItemID)
				}
				return err
			})
			if err != nil {
				slog.Error("outbox: drain", "err", err)
				continue
			}
			if sent > 0 {
				slog.Info("outbox drained", "sent", sent)
			}
		}
	}
}
