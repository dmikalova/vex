package engine

import (
	"fmt"
	"strings"
)

// OverrideStats masks the power and/or armor of every creature to a fixed value
// for the duration — The Pale Star makes each creature considered to have 1 power
// and 0 armor. It is a read-time mask installed as a continuous effect: the stored
// power counters and armor pool are untouched and revealed again when it lifts.
// The mask ignores every other modifier, so a creature's power and armor are
// exactly the masked values while it is active.
type OverrideStats struct {
	Power    int
	Armor    int
	HasPower bool
	HasArmor bool
	Duration Duration
}

// validate requires at least one masked stat and a supported duration.
func (e OverrideStats) validate() error {
	if !e.HasPower && !e.HasArmor {
		return fmt.Errorf("OverrideStats: must mask power, armor, or both")
	}
	if e.Duration != RemainderOfPlayerTurn {
		return fmt.Errorf("OverrideStats: duration must be RemainderOfPlayerTurn")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, each creature is
// considered to have 1 power and 0 armor".
func (e OverrideStats) Text() string {
	var parts []string
	if e.HasPower {
		parts = append(parts, fmt.Sprintf("%d power", e.Power))
	}
	if e.HasArmor {
		parts = append(parts, fmt.Sprintf("%d armor", e.Armor))
	}
	return "for the remainder of the turn, each creature is considered to have " +
		strings.Join(parts, " and ")
}

// Resolve installs the stat mask over every creature for the duration.
func (e OverrideStats) Resolve(ctx *EffectContext) {
	ctx.Resolver.SetStatOverride(
		StatMask{Value: int8(e.Power), Set: e.HasPower},
		StatMask{Value: int8(e.Armor), Set: e.HasArmor},
		e.Duration,
	)
}
