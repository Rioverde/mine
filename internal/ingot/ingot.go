package ingot

import (
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/ore"
)

const MaxQuality = 100

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

func FromOre(o *ore.Ore) *Ingot {
	time.Sleep(200 * time.Millisecond)
	return &Ingot{
		Ore:     o,
		Quality: rand.IntN(MaxQuality),
	}
}
