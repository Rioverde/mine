package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type KingdomRepository interface {
	GetBalance(ctx context.Context) (int64, error)
	CreateContract(ctx context.Context, contract *Contract) error
}

type pgxKingdomRepo struct {
	pool *pgxpool.Pool
}

func NewKingdomRepository(pool *pgxpool.Pool) KingdomRepository {
	return &pgxKingdomRepo{pool: pool}
}

type Contract struct {
	ID             uuid.UUID
	TargetItem     string
	TargetMaterial string
	MinQuality     int
	RequiredQty    int
	RewardCopper   int64
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Treasury struct {
	BalanceCopper int64
	UpdatedAt     time.Time
}

func (p *pgxKingdomRepo) GetBalance(ctx context.Context) (int64, error)

func (p *pgxKingdomRepo) CreateContract(ctx context.Context, contract *Contract) error
