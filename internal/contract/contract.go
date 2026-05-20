package contract

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/domain"
	"github.com/Rioverde/mine/internal/repo"
)

func GenerateContract(ctx context.Context, cfg *config.Config, r repo.KingdomRepository) {
	ticker := time.NewTicker(cfg.Economy.ContractGenerationInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			balance, err := r.GetBalance(ctx, cfg.Kingdom.ID)
			if err != nil {
				slog.Error("get balance", "err", err)
				continue
			}

			if balance <= 0 {
				slog.Info("treasury empty, skipping contract generation", "kingdom", cfg.Kingdom.ID)
				continue
			}

			c := domain.NewRandomContract(cfg.Economy, cfg.Ore, cfg.Kingdom.ID)

			if c.RewardCopper > balance {
				slog.Info("insufficient balance for contract", "balance", balance, "reward", c.RewardCopper)
				continue
			}

			if err := r.CreateContract(ctx, c); err != nil {
				slog.Error("create contract", "err", err)
				continue
			}

			slog.Info("contract created",
				"id", c.ID,
				"item", c.TargetItem,
				"qty", c.RequiredQty,
				"reward", c.RewardCopper,
			)
		}
	}

}
