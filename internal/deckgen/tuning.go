package deckgen

import "github.com/dmikalova/vactrol/internal/engine"

// Tuning holds every knob generation reads. DefaultTuning is the calibrated
// baseline; a Set starts from it and overrides only what it needs. Rates are
// per-slot Bernoulli probabilities, so their counts across the 36 slots are
// binomially distributed — the mean is the target and the spread is a natural
// bell curve, with no explicit distribution code.
type Tuning struct {
	// RarityWeights is the weighted draw of a slot's intrinsic rarity. A rarity
	// absent from the map is never rolled.
	RarityWeights map[engine.Rarity]float64

	// SpecialRate is the chance a slot is filled by a houseless Special card
	// instead of a normal rarity draw.
	SpecialRate float64
	// MaverickRate is the chance a slot's card is drawn from a different House of
	// the same Set (then rehoused to the pod's House).
	MaverickRate float64
	// LegacyRate is the chance a slot's card is drawn from the same House of a
	// different Set. It has no effect without a legacy pool (a single-set build).
	LegacyRate float64

	// DuplicateRate is the per-rarity chance a slot copies an already-placed
	// same-pod, same-rarity card instead of drawing fresh.
	DuplicateRate map[engine.Rarity]float64

	// HouseWeights biases House selection; a House absent from the map has weight
	// one. HouseExclusions lists pairs of Houses that cannot both be chosen (the
	// second is removed once the first is picked, and vice versa).
	HouseWeights    map[engine.House]float64
	HouseExclusions [][2]engine.House
}

// DefaultTuning returns the calibrated baseline: ~18 common / 12 uncommon /
// 6 rare per deck, ~3 mavericks, a special about one deck in twelve, and a couple
// of duplicates.
func DefaultTuning() Tuning {
	return Tuning{
		RarityWeights: map[engine.Rarity]float64{
			engine.Common:   0.50,
			engine.Uncommon: 0.333,
			engine.Rare:     0.167,
		},
		SpecialRate:  0.0024,
		MaverickRate: 1.0 / 12.0,
		LegacyRate:   1.0 / 6.0,
		DuplicateRate: map[engine.Rarity]float64{
			engine.Common:   0.10,
			engine.Uncommon: 0.07,
			engine.Rare:     0.04,
		},
	}
}
