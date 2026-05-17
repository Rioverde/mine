package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Rioverde/mine/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo interface {
	SaveItem(ctx context.Context, id uuid.UUID, factoryName string, it *domain.Item, payload map[string]any) error
	DrainOutbox(ctx context.Context, limit int, publish PublishFunc) (int, error)
	DeleteExpiredOutbox(ctx context.Context) (int64, error)
	Close()
}

type OutboxRow struct {
	ID      int64
	ItemID  uuid.UUID
	Payload map[string]any
}

type PublishFunc func(OutboxRow) error

type pgxRepo struct {
	pool *pgxpool.Pool
}

func Connect(ctx context.Context) (Repo, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env("STORAGE_DB_USER", "mine"),
		os.Getenv("STORAGE_DB_PASSWORD"),
		env("STORAGE_DB_HOST", "localhost"),
		env("STORAGE_DB_PORT", "5432"),
		env("STORAGE_DB_NAME", "mine"),
	)
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &pgxRepo{pool: pool}, nil
}

func (r *pgxRepo) Close() {
	r.pool.Close()
}

func (r *pgxRepo) SaveItem(ctx context.Context, id uuid.UUID, factoryName string, it *domain.Item, payload map[string]any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
        INSERT INTO items (id, name, quality, ore_material, ore_capacity, ingot_quality, factory)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		id,
		it.Name(), it.Quality,
		strings.ToLower(it.Ingot.Ore.Name()),
		it.Ingot.Ore.Capacity,
		it.Ingot.Quality,
		factoryName,
	); err != nil {
		return fmt.Errorf("insert items: %w", err)
	}

	if _, err := tx.Exec(ctx, `
        INSERT INTO outbox (item_id, payload) VALUES ($1, $2)
    `, id, payloadJSON); err != nil {
		return fmt.Errorf("insert outbox: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *pgxRepo) DrainOutbox(ctx context.Context, limit int, publish PublishFunc) (int, error) {
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

func (r *pgxRepo) DeleteExpiredOutbox(ctx context.Context) (int64, error) {
	tag, err := r.pool.Exec(ctx, `
        DELETE FROM outbox WHERE sent_at < now() - INTERVAL '7 days'
    `)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
