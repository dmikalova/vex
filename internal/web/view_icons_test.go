package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the glyph-line
// renderer in view_icons.go.

// TestSigilMarkCoversEverySigil walks sigils() and fails on any verb that renders
// an empty mark — neither a text mark nor an asset. A blank sigil draws a badge
// with nothing in it, which the card-walking collision test cannot see because two
// verbs that both draw nothing still differ in the glyph struct. It mirrors TestTargetGlyphCoversEveryTargetKind
// on the subject axis (ADR 0047).
func TestSigilMarkCoversEverySigil(t *testing.T) {
	for _, s := range sigils() {
		name := sigilName(s)
		if name == "" {
			t.Errorf("sigil %d has no CSS modifier name", s)
		}
		if sigilAsset(s) == "" && sigilText(s) == "" {
			t.Errorf("sigil %q renders no mark; give it a text mark or an asset (ADR 0047)", name)
		}
		html := app.HTMLString(sigilMark(s))
		if !strings.Contains(html, "card-glyph-sigil--"+name) {
			t.Errorf("sigil %q renders %q, want its own modifier class", name, html)
		}
	}
}

// TestSigilRendersOntoItsGlyph checks the sigil is composited inside the glyph's
// icon span — the element app.css positions it against — rather than sitting beside
// it as its own slot in the strip, which is what ADR 0047 rejected for spending a
// horizontal slot on a strip cardFitScript already squeezes.
func TestSigilRendersOntoItsGlyph(t *testing.T) {
	html := app.HTMLString(iconGlyph(glyph{asset: "aember", sigil: sigilGain}))
	host := strings.Index(html, "card-glyph-icon")
	mark := strings.Index(html, "card-glyph-sigil")
	if host < 0 || mark < host {
		t.Errorf("sigil is not nested in the glyph's icon span: %s", html)
	}
	if !strings.Contains(html, "+") {
		t.Errorf("gain sigil does not draw its mark: %s", html)
	}
}

// TestPunctMarkCoversEveryPunct walks puncts() and fails on any mark that renders
// nothing or shares a character with another. Punctuation is the strip's only
// structure channel, so two marks drawing the same character would silently fuse
// two constructs — the exact conflation ADR 0047 refused when it kept the arrow
// out of the condition grammar.
func TestPunctMarkCoversEveryPunct(t *testing.T) {
	seen := map[string]punct{}
	for _, p := range puncts() {
		name := punctName(p)
		if name == "" {
			t.Errorf("punct %d has no CSS modifier name", p)
		}
		text := punctText(p)
		if text == "" {
			t.Errorf("punct %q draws no mark (ADR 0047)", name)
			continue
		}
		if other, dup := seen[text]; dup {
			t.Errorf("punct %q draws %q, already drawn by %q", name, text, punctName(other))
		}
		seen[text] = p
		html := app.HTMLString(punctMark(p))
		if !strings.Contains(html, "card-glyph-punct--"+name) {
			t.Errorf("punct %q renders %q, want its own modifier class", name, html)
		}
	}
}

// TestPunctRendersBesideItsGlyph checks the mark is drawn before the glyph's icon
// span rather than composited onto it: punctuation sits between glyphs, where a
// sigil sits on one (ADR 0047).
func TestPunctRendersBesideItsGlyph(t *testing.T) {
	html := app.HTMLString(iconGlyph(glyph{asset: "aember", punct: punctGate}))
	mark := strings.Index(html, "card-glyph-punct")
	host := strings.Index(html, "card-glyph-icon")
	if mark < 0 || host < 0 || mark > host {
		t.Errorf("punct is not drawn before the glyph's icon span: %s", html)
	}
}

// TestComparisonMarkCoversEveryComparison walks engine.Comparisons() and fails on
// any comparison with no mark. The glyph field is the engine's own enum, so a
// comparison added to the engine reaches the strip the moment it exists; without
// this it would reach it silently, drawing a bare number that reads as an exact
// count (ADR 0047).
//
// The relative pair is the one declared overlap: MoreThanYou and MoreThanOpponent
// both mean "more than the other player's pool" and are told apart by the subject
// glyph's decor, so they share ">" rather than encoding the subject twice.
func TestComparisonMarkCoversEveryComparison(t *testing.T) {
	relative := map[engine.Comparison]bool{
		engine.MoreThanYou:      true,
		engine.MoreThanOpponent: true,
	}
	seen := map[string]engine.Comparison{}
	for _, c := range engine.Comparisons() {
		mark := comparisonMark(c)
		if mark == "" {
			t.Errorf("comparison %d draws no mark; give it one (ADR 0047)", c)
			continue
		}
		if relative[c] {
			continue
		}
		if other, dup := seen[mark]; dup {
			t.Errorf("comparison %d draws %q, already drawn by comparison %d", c, mark, other)
		}
		seen[mark] = c
	}
	if comparisonMark(comparisonNone) != "" {
		t.Errorf("an unset comparison draws %q, want no mark", comparisonMark(comparisonNone))
	}
}

// TestQtyBadgeCarriesTheComparison checks the mark rides the quantity badge rather
// than becoming a badge of its own: "power 3 or more" is one chip reading "≥3",
// not a 3 and a separate comparison mark.
func TestQtyBadgeCarriesTheComparison(t *testing.T) {
	html := app.HTMLString(iconGlyph(glyph{
		asset:      "power",
		qty:        3,
		comparison: engine.AtLeast,
	}))
	badge := strings.Index(html, "card-glyph-qty")
	mark := strings.Index(html, "card-glyph-cmp")
	if badge < 0 || mark < badge {
		t.Errorf("comparison is not nested in the quantity badge: %s", html)
	}
	if !strings.Contains(html, "≥") || !strings.Contains(html, "3") {
		t.Errorf("badge does not read as a threshold: %s", html)
	}
}

// TestParityComparisonDrawsWithoutANumber checks a comparison that ignores Amount
// still draws its badge. Even and Odd compare a quantity's shape rather than its
// size, so the mark is the whole badge, and a glyph with no qty would otherwise
// print nothing at all.
func TestParityComparisonDrawsWithoutANumber(t *testing.T) {
	html := app.HTMLString(iconGlyph(glyph{asset: "aember", comparison: engine.Even}))
	if !strings.Contains(html, "card-glyph-cmp") || !strings.Contains(html, "2n") {
		t.Errorf("a parity comparison draws no badge: %s", html)
	}
}

// TestQtyBadgeStillPrintsPlainCounts pins the behaviour the comparison rides on
// top of: a glyph with a count and no comparison prints the count alone, and a
// glyph with neither prints no badge at all.
func TestQtyBadgeStillPrintsPlainCounts(t *testing.T) {
	html := app.HTMLString(iconGlyph(glyph{asset: "damage", qty: 3}))
	if !strings.Contains(html, "card-glyph-qty") || strings.Contains(html, "card-glyph-cmp") {
		t.Errorf("a plain count draws a comparison: %s", html)
	}
	if html = app.HTMLString(
		iconGlyph(glyph{asset: "damage"}),
	); strings.Contains(
		html,
		"card-glyph-qty",
	) {
		t.Errorf("a glyph with no number draws an empty badge: %s", html)
	}
}
