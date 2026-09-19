package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/engine"
)

// iconFallbackAllowed once listed effect types the Iconography pass rendered as
// the abstract "unknown" glyph. It is gone: an unmapped mechanic is a bug, not an
// exemption. The totality tests below fail unconditionally on the fallback, so a
// new mechanic ships only once a glyph is drawn for it — the fix is always to add
// a glyph, never to allowlist the fallback (ADR 0022).

// forEachAbilityEffect walks every triggered-ability effect on every card the
// gallery shows (materialized variants included) and calls fn with it.
func forEachAbilityEffect(fn func(name string, covered bool)) {
	regs := card.Cards()
	for i := range regs {
		defs := materializedDefs(regs[i])
		for j := range defs {
			def := defs[j]
			for _, ab := range def.Abilities {
				_, covered := effectGlyphs(ab.Effect)
				fn(effectTypeName(ab.Effect), covered)
			}
		}
	}
}

// TestTargetGlyphCoversEveryTargetKind walks engine.TargetKinds() and fails on any
// kind that falls through targetGlyph's default to a text chip. The strip is pure
// icons (ADR 0022), so a worded chip is the target-side equivalent of the unknown
// glyph — and the card-walking totality tests miss it, because a text chip is not
// the unknown asset. Six context-reference kinds (TargetEachNeighbor,
// TargetEachUpgradeOnThis, TargetGrantingCard, TargetAttachedHost,
// TargetTheSameCreature, TargetFormerNeighbors) shipped that way.
func TestTargetGlyphCoversEveryTargetKind(t *testing.T) {
	for _, kind := range engine.TargetKinds() {
		g := targetGlyph(engine.Target{Kind: kind})
		if g.asset == "" {
			t.Errorf(
				"target kind %d renders as the text chip %q; give it a glyph (ADR 0022)",
				kind, g.text,
			)
		}
	}
}

// TestIconTotality walks every card and fails loud if any triggered ability's
// effect falls back to the abstract glyph instead of a real mapping. A new
// mechanic with no glyph decision fails here: the fix is to add a glyph, never to
// tolerate the fallback — there is no allowlist to exempt it (ADR 0022).
func TestIconTotality(t *testing.T) {
	forEachAbilityEffect(func(name string, covered bool) {
		if !covered {
			t.Errorf("effect %s has no glyph mapping; add one (ADR 0022)", name)
		}
	})
}

// TestNoResidualUnknownGlyph renders every card's full Icon strip and fails if any
// glyph is the abstract unknown — in a composed line's glyphs or in its trigger
// heads. Unlike TestIconTotality — which only reads the top-level covered flag —
// this walks the composed lines, so an unknown buried inside a wrapper (a
// ChooseHouseThen's inner effect, a granted ability) or at a trigger head that a
// masking covered=true would hide is still caught (ADR 0022).
func TestNoResidualUnknownGlyph(t *testing.T) {
	regs := card.Cards()
	for i := range regs {
		defs := materializedDefs(regs[i])
		for j := range defs {
			def := defs[j]
			for _, line := range cardGlyphs(&def) {
				for _, tr := range line.triggers {
					if tr == "glyph-unknown" {
						t.Errorf("%s renders an unknown trigger glyph in its Icon strip", def.Name)
					}
				}
				for _, g := range line.glyphs {
					if g.asset == "glyph-unknown" {
						t.Errorf("%s renders an unknown glyph in its Icon strip", def.Name)
					}
				}
			}
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
		Per:    engine.ForgedKeys{Player: engine.Opponent},
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
		DrawModifier: engine.DrawModifier{
			Player: engine.Controller,
			Amount: 1,
		},
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

// TestKeywordGlyphsShowsSplashAttack checks the Splash-attack combat keyword draws
// its own glyph (damage marked as reaching each neighbor) rather than nothing, so a
// creature whose only mechanic is Splash-attack still draws a strip.
func TestKeywordGlyphsShowsSplashAttack(t *testing.T) {
	def := &engine.CardDefinition{SplashAttack: 2}
	gs := keywordGlyphs(def)
	if len(gs) != 1 || gs[0].asset != "damage" || gs[0].qty != 2 ||
		gs[0].decor&decorEach == 0 {
		t.Errorf("splash glyphs = %+v, want damage x2 reaching each", gs)
	}
}

// TestStaticLinesShowsBonusAndGrant checks an upgrade whose Static grants a power
// bonus plus a triggered ability renders both a stat line and the granted ability's
// own line, so an upgrade card is not left with an empty strip.
func TestStaticLinesShowsBonusAndGrant(t *testing.T) {
	lines := staticLines(engine.StaticModifier{
		PowerBonus: 3,
		Granted: []engine.Ability{{
			Trigger: engine.TriggerAfterReap,
			Effect:  engine.Draw{Amount: 1},
		}},
	})
	if len(lines) != 2 {
		t.Fatalf("want a stat line and a granted line, got %d", len(lines))
	}
	if got := lines[0].glyphs; len(got) != 1 || got[0].asset != "power" || got[0].qty != 3 {
		t.Errorf("stat line = %+v, want power x3", got)
	}
	if got := lines[1].triggers; len(got) != 1 || got[0] != "glyph-reap" {
		t.Errorf("granted trigger = %v, want [glyph-reap]", got)
	}
}

// TestRestrictionLinesBansReaping checks a card that stops the opponent's creatures
// from reaping renders the reap glyph, tinted to the enemy, struck by the ban glyph.
func TestRestrictionLinesBansReaping(t *testing.T) {
	lines := restrictionLines(engine.Restrictions{Reaping: engine.Opponent})
	if len(lines) != 1 {
		t.Fatalf("want one restriction line, got %d", len(lines))
	}
	gs := lines[0].glyphs
	if len(gs) != 2 || gs[0].asset != "glyph-reap" ||
		gs[0].decor&decorEnemy == 0 || gs[1].asset != "glyph-ban" {
		t.Errorf("reaping ban glyphs = %+v, want enemy glyph-reap + ban", gs)
	}
}

// TestRestrictionLinesShowsToll checks a card that charges the opponent Æmber to
// take an artifact action (a toll) renders the artifact/action glyph beside the
// enemy Æmber owed.
func TestRestrictionLinesShowsToll(t *testing.T) {
	lines := restrictionLines(engine.Restrictions{
		Toll: engine.Toll{
			Action: engine.TollUseArtifact,
			Amount: 1,
		},
	})
	if len(lines) != 1 {
		t.Fatalf("want one toll line, got %d", len(lines))
	}
	gs := lines[0].glyphs
	if len(gs) != 2 || gs[0].asset != "glyph-action" ||
		gs[1].asset != "aember" || gs[1].qty != 1 || gs[1].decor&decorEnemy == 0 {
		t.Errorf("toll glyphs = %+v, want glyph-action + enemy aember x1", gs)
	}
}

// TestCardFeatureLinesShowsAemberCannotBeStolen checks a card whose only mechanic is
// the card-level AemberCannotBeStolen flag draws a strip (enemy Æmber struck by the
// ban glyph) rather than none.
func TestCardFeatureLinesShowsAemberCannotBeStolen(t *testing.T) {
	lines := cardFeatureLines(&engine.CardDefinition{AemberCannotBeStolen: engine.AlwaysMet{}})
	if len(lines) != 1 {
		t.Fatalf("want one feature line, got %d", len(lines))
	}
	gs := lines[0].glyphs
	if len(gs) != 2 || gs[0].asset != "aember" ||
		gs[0].decor&decorEnemy == 0 || gs[1].asset != "glyph-ban" {
		t.Errorf("feature glyphs = %+v, want enemy aember + ban", gs)
	}
}

// TestReplaceGlyphsCapture checks a card that replaces an Æmber-flow event with a
// capture (Ether Spider) renders the affected pool's Æmber swapped for enemy Æmber.
func TestReplaceGlyphsCapture(t *testing.T) {
	gs := replaceGlyphs(engine.Instead{
		Of:     engine.EventAemberAddedToPool,
		Player: engine.Opponent,
		With:   engine.Capture,
	})
	if len(gs) != 3 || gs[0].asset != "aember" ||
		gs[1].asset != "glyph-swap" || gs[2].asset != "aember" ||
		gs[2].decor&decorEnemy == 0 {
		t.Errorf("replace glyphs = %+v, want aember swap→ enemy aember", gs)
	}
}

// TestTriggerIconPhase checks the phase triggers (start/end of turn, end of ready
// step) share the turn-phase glyph rather than falling back to the unknown glyph.
func TestTriggerIconPhase(t *testing.T) {
	for _, tr := range []engine.Trigger{
		engine.TriggerStartOfTurn, engine.TriggerEndOfTurn, engine.TriggerEndOfReadyStep,
	} {
		if got := triggerIcon(tr); got != "phase-turn" {
			t.Errorf("triggerIcon(%v) = %q, want phase-turn", tr, got)
		}
	}
}
