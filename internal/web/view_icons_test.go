package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is part of the Iconography pass (ADR 0022): tests for the glyph-line
// renderer in view_icons.go.

// TestSigilMarkCoversEverySigil walks sigils() and fails on any verb that renders
// an empty mark — neither a text mark nor an asset. A blank sigil is the verb-side
// equivalent of the unknown glyph: it draws a badge with nothing in it, which the
// card-walking collision test cannot see because two verbs that both draw nothing
// still differ in the glyph struct. It mirrors TestTargetGlyphCoversEveryTargetKind
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
