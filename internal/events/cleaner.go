package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/repo"
)

func OutboxCleaner(ctx context.Context, r repo.Repo, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			deleted, err := r.DeleteExpiredOutbox(ctx)
			if err != nil {
				slog.Error("outbox cleaner", "err", err)
				continue
			}
			if deleted > 0 {
				slog.Info("outbox cleaned", "deleted", deleted)
			}
		}
	}
}
