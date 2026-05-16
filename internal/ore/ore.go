package ore

import (
	"math/rand/v2"

	"github.com/Rioverde/mine/internal/config"
)

const (
	Iron = iota
	Copper
	Gold
)

type Ore struct {
	Material int
	Capacity int
}

func (o *Ore) Name() string {
	switch o.Material {
	case Iron:
		return "Iron"
	case Copper:
		return "Copper"
	case Gold:
		return "Gold"
	}
	return "Unknown"
}

func Random(cfg config.OreConfig) *Ore {
	return &Ore{
		Material: rand.IntN(cfg.NumberOfOres),
		Capacity: rand.IntN(cfg.MaxCapacity),
	}
}
