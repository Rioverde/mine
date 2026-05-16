package ore

import "math/rand/v2"

const (
	Iron = iota
	Copper
	Gold
)

const (
	MaxCapacity  = 100
	NumberOfOres = 3
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

func Random() *Ore {
	return &Ore{
		Material: rand.IntN(NumberOfOres),
		Capacity: rand.IntN(MaxCapacity),
	}
}
