package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestDamageEffectGlyphsTranscription binds the damage family to a few worked
// examples so the grammar cannot drift silently.
func TestDamageEffectGlyphsTranscription(t *testing.T) {
	// "Deal 3 damage to an enemy creature" → damage(3) → creature[enemy, chosen].
	gs, covered := effectGlyphs(engine.DealDamage{
		Amount: 3,
		Target: engine.Target{Kind: engine.TargetChosenEnemyCreature},
	})
	if !covered {
		t.Fatal("DealDamage to a chosen enemy creature should be covered")
	}
	if len(gs) != 2 {
		t.Fatalf("want 2 glyphs, got %d", len(gs))
	}
	if gs[0].asset != "damage" || gs[0].qty != 3 {
		t.Errorf("first glyph = %+v, want damage qty 3", gs[0])
	}
	if gs[1].asset != "type-creature" || gs[1].punct != punctArrow ||
		gs[1].decor&decorEnemy == 0 || gs[1].decor&decorChosen == 0 {
		t.Errorf("second glyph = %+v, want arrowed enemy chosen creature", gs[1])
	}

	// "Destroy a creature" → destroy → creature[chosen].
	gs, covered = effectGlyphs(engine.Destroy{
		Target: engine.Target{Kind: engine.TargetChosenCreature},
	})
	if !covered || len(gs) != 2 || gs[0].asset != "glyph-destroy" ||
		gs[1].asset != "type-creature" || gs[1].punct != punctArrow {
		t.Errorf("Destroy glyphs = %+v (covered=%v), want destroy → chosen creature", gs, covered)
	}
}

// TestDelayedWipeDrawsItsWindow pins the pair the duration channel exists to
// separate: engine.BatchDestroy over every creature and
// engine.DestroyEachCreatureAtEndOfTurn are the same wipe, and the only difference
// is that one is scheduled into the end-of-turn window. They drew byte-identical
// strips until the turn marker was appended (ADR 0047).
func TestDelayedWipeDrawsItsWindow(t *testing.T) {
	delayed, covered := effectGlyphs(engine.DestroyEachCreatureAtEndOfTurn{})
	if !covered {
		t.Fatal("the end-of-turn wipe should be covered")
	}
	if len(delayed) == 0 || delayed[len(delayed)-1].asset != "phase-turn" {
		t.Fatalf("glyphs = %+v, want the turn marker closing the run", delayed)
	}
	now, _ := effectGlyphs(engine.Destroy{
		Target: engine.Target{Kind: engine.TargetEachCreature},
	})
	if fmt.Sprintf("%+v", now) == fmt.Sprintf("%+v", delayed) {
		t.Errorf("the immediate and delayed wipes still draw the same strip %+v", now)
	}
}
