package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the shared
// noun and decor vocabulary in icon_target.go.

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
