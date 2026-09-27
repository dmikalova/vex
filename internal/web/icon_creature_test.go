package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the
// creature-state family in icon_creature.go.

// TestCreatureStateDrawsItsOwnArt pins the three conditions to the assets drawn for
// them, which the board already uses. Ward drew shield.svg — armour's icon — so
// Ward, "cannot be dealt damage" and armour were one picture, and Enrage drew the
// fight glyph while enrage.svg sat unused in web/assets (ADR 0047).
func TestCreatureStateDrawsItsOwnArt(t *testing.T) {
	for _, tc := range []struct {
		name  string
		e     engine.Effect
		asset string
	}{
		{"ward", engine.Ward{Target: engine.Target{Kind: engine.TargetChosenCreature}}, "ward"},
		{"enrage", engine.Enrage{Target: engine.Target{Kind: engine.TargetChosenCreature}}, "enrage"},
		{"stun", engine.Stun{Target: engine.Target{Kind: engine.TargetChosenCreature}}, "stun"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gs, covered := effectGlyphs(tc.e)
			if !covered || len(gs) == 0 || gs[0].asset != tc.asset {
				t.Errorf(
					"%s glyphs = %+v (covered=%v), want %q first",
					tc.name,
					gs,
					covered,
					tc.asset,
				)
			}
		})
	}
}

// TestCountersReadAsPlacedOrRemoved checks that placing and removing a counter draw
// the same counter noun under opposite verb sigils. They were byte-identical before
// the verb axis, so a player could not tell one from the other (ADR 0047).
func TestCountersReadAsPlacedOrRemoved(t *testing.T) {
	target := engine.Target{Kind: engine.TargetChosenCreature}
	place, _ := effectGlyphs(engine.PlaceCounter{Kind: engine.CounterDoom, Target: target})
	remove, _ := effectGlyphs(engine.RemoveCounters{Kind: engine.CounterDoom, Target: target})
	if place[0].asset != remove[0].asset {
		t.Errorf(
			"place %q and remove %q should share the counter noun",
			place[0].asset,
			remove[0].asset,
		)
	}
	if place[0].sigil != sigilGain || remove[0].sigil != sigilLose {
		t.Errorf(
			"place sigil = %v, remove sigil = %v, want gain and lose",
			place[0].sigil, remove[0].sigil,
		)
	}
}
