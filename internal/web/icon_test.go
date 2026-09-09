package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/engine"
)

// iconFallbackAllowed lists the effect types the Iconography pass does not yet map
// to a meaningful glyph and so renders as the abstract "unknown" glyph. It is a
// shrinking triage list, not a permanent exemption: as glyphs are drawn for these
// mechanics their names come off. The totality test below fails when a card uses
// an effect that is neither mapped nor on this list, so a new mechanic cannot ship
// without a glyph decision (ADR 0022). The list is currently empty: every effect a
// card uses maps to a glyph. A new unmapped mechanic either gets a glyph here or an
// entry added to this list with a reason.
var iconFallbackAllowed = map[string]bool{}

// forEachAbilityEffect walks every triggered-ability effect on every card the
// gallery shows (materialized variants included) and calls fn with it.
func forEachAbilityEffect(fn func(name string, covered bool)) {
	regs := card.Cards()
	for i := range regs {
		for _, def := range materializedDefs(regs[i]) {
			def := def
			for _, ab := range def.Abilities {
				_, covered := effectGlyphs(ab.Effect)
				fn(effectTypeName(ab.Effect), covered)
			}
		}
	}
}

// TestIconTotality walks every card and fails loud if any triggered ability's
// effect is neither mapped to a glyph nor on the shrinking fallback allowlist. A
// new mechanic with no glyph decision fails here (ADR 0022).
func TestIconTotality(t *testing.T) {
	forEachAbilityEffect(func(name string, covered bool) {
		if !covered && !iconFallbackAllowed[name] {
			t.Errorf("effect %s has no glyph mapping and is not on iconFallbackAllowed", name)
		}
	})
}

// TestIconFallbackAllowlistIsUsed keeps the allowlist honest: an entry that no
// card exercises any more should come off, so the list shrinks toward empty as
// glyphs are drawn.
func TestIconFallbackAllowlistIsUsed(t *testing.T) {
	used := map[string]bool{}
	forEachAbilityEffect(func(name string, covered bool) {
		if !covered {
			used[name] = true
		}
	})
	for name := range iconFallbackAllowed {
		if !used[name] {
			t.Errorf("iconFallbackAllowed lists %s, but no card falls back to it; remove it", name)
		}
	}
}

// TestEffectGlyphsTranscription binds the pass to a few worked examples so the
// grammar cannot drift silently.
func TestEffectGlyphsTranscription(t *testing.T) {
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
	if gs[1].asset != "type-creature" || !gs[1].arrow ||
		gs[1].decor&decorEnemy == 0 || gs[1].decor&decorChosen == 0 {
		t.Errorf("second glyph = %+v, want arrowed enemy chosen creature", gs[1])
	}

	// "Destroy a creature" → destroy → creature[chosen].
	gs, covered = effectGlyphs(engine.Destroy{
		Target: engine.Target{Kind: engine.TargetChosenCreature},
	})
	if !covered || len(gs) != 2 || gs[0].asset != "glyph-destroy" ||
		gs[1].asset != "type-creature" || !gs[1].arrow {
		t.Errorf("Destroy glyphs = %+v (covered=%v), want destroy → chosen creature", gs, covered)
	}

	// A Per variant of GainAember maps to the Æmber glyph without a numeral: the
	// count scales by a board quantity that lives in the text, not on the glyph.
	gs, covered = effectGlyphs(engine.GainAember{
		Player: engine.Controller,
		Amount: 1,
		Per:    engine.OpponentForgedKeys{},
	})
	if !covered || len(gs) != 1 || gs[0].asset != "aember" || gs[0].qty != 0 {
		t.Errorf(
			"GainAember Per glyphs = %+v (covered=%v), want a single æmber glyph with no numeral",
			gs,
			covered,
		)
	}
}

// TestCardGlyphsMergesActionTriggers checks that a Play/Fight/Reap ability — three
// abilities sharing one effect on the three action triggers — renders as one line
// carrying all three trigger glyphs, the icon counterpart of the "Play/Fight/Reap:"
// rules-text shorthand, rather than three repeated lines.
func TestCardGlyphsMergesActionTriggers(t *testing.T) {
	effect := engine.Destroy{Target: engine.Target{Kind: engine.TargetChosenCreature}}
	def := &engine.CardDefinition{Abilities: []engine.Ability{
		{Trigger: engine.TriggerAfterPlay, Effect: effect},
		{Trigger: engine.TriggerAfterFight, Effect: effect},
		{Trigger: engine.TriggerAfterReap, Effect: effect},
	}}
	lines := cardGlyphs(def)
	if len(lines) != 1 {
		t.Fatalf("want one merged line, got %d", len(lines))
	}
	if got := lines[0].triggers; len(got) != 3 ||
		got[0] != "glyph-play" || got[1] != "glyph-fight" || got[2] != "glyph-reap" {
		t.Errorf("triggers = %v, want [glyph-play glyph-fight glyph-reap]", got)
	}
}

// TestCardGlyphsShowsDrawModifier checks a card whose only mechanic is a continuous
// hand-refill change (Mother: draw +1) still draws a glyph strip rather than none.
func TestCardGlyphsShowsDrawModifier(t *testing.T) {
	def := &engine.CardDefinition{
		DrawModifier: engine.DrawModifier{Player: engine.Controller, Amount: 1},
	}
	lines := cardGlyphs(def)
	if len(lines) != 1 {
		t.Fatalf("want one draw-modifier line, got %d", len(lines))
	}
	if got := lines[0].glyphs; len(got) != 1 ||
		got[0].asset != "zone-hand" || got[0].qty != 1 ||
		got[0].decor&decorFriendly == 0 {
		t.Errorf("draw-modifier glyphs = %+v, want a friendly zone-hand +1", got)
	}
}
