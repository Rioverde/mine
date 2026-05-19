package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OutboxRepository interface {
	DrainOutbox(ctx context.Context, limit int, publish PublishFunc) (int, error)
	DeleteExpiredOutbox(ctx context.Context) (int64, error)
}

type OutboxRow struct {
	ID      int64
	ItemID  uuid.UUID
	Payload map[string]any
}

// PublishFunc is called for each fetched outbox row. Returning nil marks the row as sent.
// Returning an error stops draining; already-published rows in the batch are still committed.
type PublishFunc func(OutboxRow) error

type pgOutboxRepo struct {
	pool *pgxpool.Pool
}

// DrainOutbox claims up to `limit` pending rows via FOR UPDATE SKIP LOCKED,
// invokes publish() for each, marks successful ones as sent, and commits atomically.
// Safe for concurrent publishers — each claims a disjoint subset.
func (r *pgOutboxRepo) DrainOutbox(ctx context.Context, limit int, publish PublishFunc) (int, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
        SELECT id, item_id, payload
        FROM outbox
        WHERE sent_at IS NULL
        ORDER BY id ASC
        LIMIT $1
        FOR UPDATE SKIP LOCKED
    `, limit)
	if err != nil {
		return 0, fmt.Errorf("select outbox: %w", err)
	}

	batch := make([]OutboxRow, 0, limit)
	for rows.Next() {
		var row OutboxRow
		var payloadJSON []byte
		if err := rows.Scan(&row.ID, &row.ItemID, &payloadJSON); err != nil {
			rows.Close()
			return 0, fmt.Errorf("scan outbox: %w", err)
		}
		if err := json.Unmarshal(payloadJSON, &row.Payload); err != nil {
			rows.Close()
			return 0, fmt.Errorf("unmarshal payload id=%d: %w", row.ID, err)
		}
		batch = append(batch, row)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	if len(batch) == 0 {
		return 0, tx.Commit(ctx)
	}

	sentIDs := make([]int64, 0, len(batch))
	for _, row := range batch {
		if err := publish(row); err != nil {
			break
		}
		sentIDs = append(sentIDs, row.ID)
	}

	if len(sentIDs) > 0 {
		if _, err := tx.Exec(ctx, `
            UPDATE outbox SET sent_at = now() WHERE id = ANY($1)
        `, sentIDs); err != nil {
			return 0, fmt.Errorf("mark sent: %w", err)
		}
	}

	return len(sentIDs), tx.Commit(ctx)
}

func (r *pgOutboxRepo) DeleteExpiredOutbox(ctx context.Context) (int64, error) {
	tag, err := r.pool.Exec(ctx, `
        DELETE FROM outbox WHERE sent_at < now() - INTERVAL '7 days'
    `)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
