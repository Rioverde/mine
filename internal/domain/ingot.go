package domain

import (
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/config"
)

type Ingot struct {
	Ore     *Ore
	Quality int
}

func (i *Ingot) Name() string {
	switch i.Ore.Material {
	case Iron:
		return "Iron Ingot"
	case Copper:
		return "Copper Ingot"
	case Gold:
		return "Gold Ingot"
	}
	return "Unknown"
}

func NewSmelter(cfg config.IngotConfig) func(*Ore) *Ingot {
	return func(o *Ore) *Ingot {
		time.Sleep(200 * time.Millisecond)
		return &Ingot{
			Ore:     o,
			Quality: rand.IntN(cfg.MaxQuality),
		}
	}
}
