package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the
// continuous-rules renderers in icon_static.go.

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
