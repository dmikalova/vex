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

// TestEveryGatherDraws holds the family's gate from the other side: every
// catalogued gatherer draws a real strip of existing assets rather than the
// abstract unknown glyph, so a batch that takes one creature from each of two
// pools cannot draw what a board wipe draws.
func TestEveryGatherDraws(t *testing.T) {
	assertFamilyDraws(t, "Gather")
}

// TestBatchDestroyDrawsItsGatherer pins what the gatherer is for: the two
// gatherers take different creatures, so the one verb carrying them must draw
// them apart. BatchDestroy drew a fixed "each creature" for both until the
// gatherer was drawn, which was a wipe's picture on Imp-losion's two picks.
func TestBatchDestroyDrawsItsGatherer(t *testing.T) {
	picks, covered := effectGlyphs(engine.BatchDestroy{Gather: engine.ChosenFromEach{
		engine.Target{Kind: engine.TargetChosenFriendlyCreature},
		engine.Target{Kind: engine.TargetChosenEnemyCreature},
	}})
	if !covered {
		t.Fatal("a batch destroy over two pools should be covered")
	}
	if len(picks) != 3 {
		t.Fatalf("two-pool batch = %+v, want the verb and both pools", picks)
	}
	if picks[1].decor&decorFriendly == 0 || picks[2].decor&decorEnemy == 0 {
		t.Errorf("two-pool batch = %+v, want a friendly pool and an enemy one", picks)
	}
	unless, _ := effectGlyphs(engine.BatchDestroy{Gather: engine.EachPlayerUnless{
		Spare: engine.CardsInPlay{Player: engine.Controller, Type: engine.Creature},
		Take:  engine.MostPowerful,
	}})
	if fmt.Sprintf("%+v", picks) == fmt.Sprintf("%+v", unless) {
		t.Errorf("both gatherers draw the same strip %+v", picks)
	}
}

// TestEachPlayerUnlessGatesItsSpareBoard checks that the board sparing a player
// is drawn behind the gate mark and closed by the refusal mark, which is what
// makes it read as "unless". Without the refusal mark the strip would say the
// opposite of the rule: that the board is what gets a player taken.
func TestEachPlayerUnlessGatesItsSpareBoard(t *testing.T) {
	gs, covered := gatherGlyphs(engine.EachPlayerUnless{
		Spare: engine.CardsInPlay{Player: engine.Controller, Type: engine.Creature},
		Take:  engine.MostPowerful,
	})
	if !covered || len(gs) < 3 {
		t.Fatalf("unless-gatherer = %+v (covered=%v), want a gated clause", gs, covered)
	}
	if gs[1].punct != punctGate {
		t.Errorf("the spare board opens with %v, want the gate mark", gs[1].punct)
	}
	if gs[len(gs)-1].asset != "glyph-ban" {
		t.Errorf("unless-gatherer = %+v, want the refusal mark closing it", gs)
	}
}

// TestEverySpreadDraws holds the Spread family's gate from the other side: every
// catalogued spread draws a real strip of existing assets, and no two of them
// draw the same one. The second half is the point of this family — all five
// spreads pick creatures, and they drew one flattened chosen-creature noun until
// each was given the shape of its own pick (ADR 0047).
func TestEverySpreadDraws(t *testing.T) {
	assertFamilyDraws(t, "Spread")
	assertFamilyRowsDrawApart(t, "Spread")
}

// TestNeighborScopesDrawApart pins the mechanic inside CreatureAndNeighbors:
// splashing every neighbor (Booby Trap) and splashing one neighbor the controller
// chooses (Mighty Lance) are different rules, so the stack and the chosen outline
// must tell them apart rather than both drawing "the neighbors".
func TestNeighborScopesDrawApart(t *testing.T) {
	seen := map[decor]engine.NeighborScope{}
	for _, scope := range engine.NeighborScopes() {
		gs, covered := spreadGlyphs(engine.CreatureAndNeighbors{
			Amount: 2,
			Splash: 1,
			Scope:  scope,
		})
		if !covered || len(gs) != 3 {
			t.Fatalf("scope %v drew %+v (covered=%v), want three glyphs", scope, gs, covered)
		}
		if prev, dup := seen[gs[2].decor]; dup {
			t.Errorf("neighbor scopes %v and %v both draw decor %v", prev, scope, gs[2].decor)
		}
		seen[gs[2].decor] = scope
	}
}

// TestEveryPerTargetDraws holds the PerTarget family's gate from the other side:
// every axis an amount can vary along draws a real, distinct strip of existing
// assets, so "for each Æmber on it" cannot draw what "for each upgrade on it"
// draws.
func TestEveryPerTargetDraws(t *testing.T) {
	assertFamilyDraws(t, "PerTarget")
	assertFamilyRowsDrawApart(t, "PerTarget")
}

// TestPerTargetDamageDrawsWhatItMeasures checks the axis reaches the strip of the
// verb carrying it: Word of Returning deals damage to each enemy creature for
// each Æmber on that creature, and without the axis drawn its strip said only
// "damage to each enemy creature" — the same picture a flat sweep draws.
func TestPerTargetDamageDrawsWhatItMeasures(t *testing.T) {
	gs, covered := effectGlyphs(engine.DealDamage{
		Amount:    1,
		Target:    engine.Target{Kind: engine.TargetEachEnemyCreature},
		PerTarget: engine.AemberOnIt,
	})
	if !covered || len(gs) != 3 {
		t.Fatalf("per-target damage = %+v (covered=%v), want three glyphs", gs, covered)
	}
	if gs[2].asset != "aember" || gs[2].decor&decorChosen == 0 {
		t.Errorf("per-target damage closes with %+v, want the Æmber it measures", gs[2])
	}
	flat, _ := effectGlyphs(engine.DealDamage{
		Amount: 1,
		Target: engine.Target{Kind: engine.TargetEachEnemyCreature},
	})
	if fmt.Sprintf("%+v", flat) == fmt.Sprintf("%+v", gs) {
		t.Errorf("a per-target sweep still draws the flat sweep's strip %+v", flat)
	}
}
