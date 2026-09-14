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
	// InterloperRate is the per-pod chance a whole House pod is drawn from the
	// legacy pool — every slot a same-House card from another Set — rather than
	// from this Set's own pool. It is the House-level counterpart of the slot-level
	// LegacyRate, and like it has no effect without a legacy pool. Rolled once per
	// pod, so a deck may carry more than one interloper pod.
	InterloperRate float64
	// ErrantRate is the per-pod chance a pod's House is replaced by a foreign House
	// — one not native to this Set — and the whole pod drawn from the legacy pool
	// as that foreign House. It is the errant counterpart of InterloperRate, which
	// keeps a native House; like it, it needs a legacy pool, and it also needs that
	// pool to carry a foreign House (else the roll never fires). Rolled once per
	// pod, before the interloper roll, so a deck may carry more than one errant pod.
	ErrantRate float64

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
// 6 rare per deck, ~3 mavericks, a special about one deck in twelve, and no
// duplicates.
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
		// InterloperRate is per pod; across a deck's three pods a rate of 1/36
		// lands an interloper pod in about one deck in twelve.
		InterloperRate: 1.0 / 36.0,
		// ErrantRate is per pod; across a deck's three pods a rate of 1/180 lands
		// an errant pod in about one deck in sixty — rarer than an interloper, as a
		// foreign House is a bigger departure than a legacy-drawn native pod.
		ErrantRate: 1.0 / 180.0,
		// Duplicate-pull is off by default: with the implemented card pool still
		// small, real decks came out with noticeably more repeats than they should.
		// The mechanic stays wired up — raise these rates on a Set to turn it back on.
		DuplicateRate: map[engine.Rarity]float64{
			engine.Common:   0,
			engine.Uncommon: 0,
			engine.Rare:     0,
		},
	}
}
