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
			continue
		}
		if !assetExists(t, g.asset) {
			t.Errorf("target kind %d: no web/assets/%s.svg", kind, g.asset)
		}
	}
}

// TestDestinationGlyphCoversEveryDestination walks engine.Destinations() and
// fails on a destination that draws no zone. destinationGlyph returns "" for an
// unmapped one and the zone-move effects reading it fall back to the unknown
// glyph, so the hole only shows once a card names that destination; the census
// guards Destinations() complete against the source, so this catches it first
// (ADR 0022). The Yours variants redirect the same zone, so they draw the same
// glyph and are covered by the plain ones.
func TestDestinationGlyphCoversEveryDestination(t *testing.T) {
	for _, d := range engine.Destinations() {
		name := destinationGlyph(d)
		if name == "" {
			t.Errorf("destination %v draws no zone glyph; give it one (ADR 0022)", d)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("destination %v: no web/assets/%s.svg", d, name)
		}
	}
}

// TestDeckDestZoneCoversEveryDeckDest walks engine.DeckDests() and fails on a
// deck routing that draws no zone. An empty asset is dropped silently at render,
// so a card that looks at the top of its deck and puts the leftovers somewhere
// undrawn simply loses that half of its strip — IntoBottomOfDeck shipped that way
// (ADR 0022).
func TestDeckDestZoneCoversEveryDeckDest(t *testing.T) {
	for _, d := range engine.DeckDests() {
		name := deckDestZone(d)
		if name == "" {
			t.Errorf("deck routing %v draws no zone glyph; give it one (ADR 0022)", d)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("deck routing %v: no web/assets/%s.svg", d, name)
		}
	}
}

// TestVerbGlyphsCoverEveryCreatureVerb walks engine.CreatureVerbCatalog() and
// fails on a verb verbGlyphs renders as nothing. The type switch there has no
// default arm, so an unlisted verb contributes no glyph at all and the strip
// quietly says the creature is chosen and nothing happens to it — the census is
// what makes a new verb reach this test before a card does (ADR 0022).
func TestVerbGlyphsCoverEveryCreatureVerb(t *testing.T) {
	for _, row := range engine.CreatureVerbCatalog() {
		gs := verbGlyphs([]engine.CreatureVerb{row.Node})
		if len(gs) != 1 || gs[0].asset == "" {
			t.Errorf("creature verb %T draws %+v, want one glyph (ADR 0022)", row.Node, gs)
			continue
		}
		if !assetExists(t, gs[0].asset) {
			t.Errorf("creature verb %T: no web/assets/%s.svg", row.Node, gs[0].asset)
		}
	}
}

// TestCounterIconNamesHaveAssets guards the generic-counter icon seam: every
// counter kind the engine enumerates resolves to a unique icon asset that exists
// on disk. A new CounterKind with no counterAsset entry (or one reusing another
// kind's icon) fails here, so the web strip can never draw two counters the same
// or fall back to a missing SVG.
func TestCounterIconNamesHaveAssets(t *testing.T) {
	if got := counterAsset(engine.CounterNone); got != "" {
		t.Errorf("counterAsset(CounterNone) = %q, want empty", got)
	}
	seen := map[string]engine.CounterKind{}
	for _, k := range engine.CounterKinds() {
		name := counterAsset(k)
		if name == "" {
			t.Errorf("counter kind %v has no icon asset", k)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("counter kind %v: no web/assets/%s.svg", k, name)
		}
		if prev, dup := seen[name]; dup {
			t.Errorf(
				"counter kinds %v and %v share icon %q; each needs a unique one",
				prev,
				k,
				name,
			)
		}
		seen[name] = k
	}
}
