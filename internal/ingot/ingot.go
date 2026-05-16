package ingot

import (
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/ore"
)

type Ingot struct {
	Ore     *ore.Ore
	Quality int
}

func (i *Ingot) Name() string {
	switch i.Ore.Material {
	case ore.Iron:
		return "Iron Ingot"
	case ore.Copper:
		return "Copper Ingot"
	case ore.Gold:
		return "Gold Ingot"
	}
	return "Unknown"
}

func NewSmelter(cfg config.IngotConfig) func(*ore.Ore) *Ingot {
	return func(o *ore.Ore) *Ingot {
		time.Sleep(200 * time.Millisecond)
		return &Ingot{
			Ore:     o,
			Quality: rand.IntN(cfg.MaxQuality),
		}
	}
}
