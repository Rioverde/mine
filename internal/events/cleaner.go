package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/repo"
)

// OutboxCleaner periodically removes outbox rows older than the retention window
// (currently hardcoded to 7 days in repo.DeleteExpiredOutbox).
// Stops only when ctx is cancelled.
func OutboxCleaner(ctx context.Context, r repo.OutboxRepository, interval time.Duration) {
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
