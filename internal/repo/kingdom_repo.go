package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rioverde/mine/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrKingdomNotFound = errors.New("kingdom not found")

type KingdomRepository interface {
	EnsureKingdom(ctx context.Context, id uuid.UUID, name string) error
	GetBalance(ctx context.Context, kingdomID uuid.UUID) (int64, error)
	CreateContract(ctx context.Context, contract *domain.Contract) error
}

type pgxKingdomRepo struct {
	pool *pgxpool.Pool
}

// EnsureKingdom upserts the kingdom row and its treasury so callers can
// assume both exist. Idempotent — safe on every startup.
func (p *pgxKingdomRepo) EnsureKingdom(ctx context.Context, id uuid.UUID, name string) error {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO kingdoms (id, name) VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING
	`, id, name); err != nil {
		return fmt.Errorf("insert kingdom: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO treasury (kingdom_id) VALUES ($1)
		ON CONFLICT (kingdom_id) DO NOTHING
	`, id); err != nil {
		return fmt.Errorf("insert treasury: %w", err)
	}

	return tx.Commit(ctx)
}

func (p *pgxKingdomRepo) GetBalance(ctx context.Context, kingdomID uuid.UUID) (int64, error) {
	var balance int64
	err := p.pool.QueryRow(ctx, `
		SELECT balance_copper FROM treasury WHERE kingdom_id = $1
	`, kingdomID).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrKingdomNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("select treasury: %w", err)
	}
	return balance, nil
}

func (p *pgxKingdomRepo) CreateContract(ctx context.Context, c *domain.Contract) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO contracts
			(id, kingdom_id, target_item, target_material, min_quality, required_qty, reward_copper, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, c.ID, c.KingdomID, c.TargetItem, c.TargetMaterial, c.MinQuality, c.RequiredQty, c.RewardCopper, c.Status)
	if err != nil {
		return fmt.Errorf("insert contract: %w", err)
	}
	return nil
}
