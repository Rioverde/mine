package domain

import (
	"errors"
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/google/uuid"
)

type ContractStatus string

const (
	ContractActive    ContractStatus = "active"
	ContractFulfilled ContractStatus = "fulfilled"
	ContractExpired   ContractStatus = "expired"
)

type Contract struct {
	ID             uuid.UUID
	KingdomID      uuid.UUID
	TargetItem     string
	TargetMaterial int
	MinQuality     int
	RequiredQty    int
	RewardCopper   int64
	Status         ContractStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewContract(kingdomID uuid.UUID, targetItem string, material int, minQ, qty int, reward int64) (*Contract, error) {
	if qty <= 0 {
		return nil, errors.New("required_qty must be positive")
	}
	if reward <= 0 {
		return nil, errors.New("reward must be positive")
	}
	if minQ < 0 || minQ > 100 {
		return nil, errors.New("min_quality out of range")
	}
	now := time.Now().UTC()
	return &Contract{
		ID:             uuid.New(),
		KingdomID:      kingdomID,
		TargetItem:     targetItem,
		TargetMaterial: material,
		MinQuality:     minQ,
		RequiredQty:    qty,
		RewardCopper:   reward,
		Status:         ContractActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (c *Contract) Matches(it *Item) bool {
	return c.Status == ContractActive &&
		it.Name() == c.TargetItem &&
		it.Ingot.Ore.Material == c.TargetMaterial &&
		it.Quality >= c.MinQuality
}

func itemNameFor(material int) string {
	switch material {
	case Iron:
		return "Sword"
	case Copper:
		return "Shield"
	case Gold:
		return "Axe"
	}
	return "Unknown"
}

func NewRandomContract(eco config.EconomyConfig, ore config.OreConfig, kingdomID uuid.UUID) *Contract {
	material := rand.IntN(ore.NumberOfOres)
	rewardSpan := eco.MaxContractReward - eco.MinContractReward + 1
	now := time.Now().UTC()
	return &Contract{
		ID:             uuid.New(),
		KingdomID:      kingdomID,
		TargetItem:     itemNameFor(material),
		TargetMaterial: material,
		MinQuality:     rand.IntN(50),
		RequiredQty:    1 + rand.IntN(5),
		RewardCopper:   eco.MinContractReward + rand.Int64N(rewardSpan),
		Status:         ContractActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
