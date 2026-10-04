package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
)

// These tests pin the player-facing glyph legend on /rulebook: that it is built
// from the census rather than a hand-kept table, and that every term it has
// glyphs for reaches the page.

// TestTermGlyphsCoverEveryDrawnTerm checks the legend is complete from the other
// side: every rule-bearing row of a gated family draws a glyph (the icon gate), so
// every term such a row names must appear in the legend. A term that fell out
// would be a picture a player meets on a card and cannot look up.
func TestTermGlyphsCoverEveryDrawnTerm(t *testing.T) {
	legend := termGlyphs()
	forEachGatedIconRow(func(_ iconFamily, row iconRow) {
		if len(legend[row.rules.Term]) == 0 {
			t.Errorf("%s owes the term %q, which the legend has no glyph for",
				row.name, row.rules.Term)
		}
	})
}

// TestTermGlyphsNameRealRuleTerms checks every key is a rulebook title the engine
// registers, so a legend entry cannot be filed under a heading /rulebook never
// prints.
func TestTermGlyphsNameRealRuleTerms(t *testing.T) {
	titles := map[string]bool{}
	for _, term := range engine.RuleTerms() {
		titles[term.Title] = true
	}
	for term := range termGlyphs() {
		if !titles[term] {
			t.Errorf("the legend files glyphs under %q, which is no rulebook term", term)
		}
	}
}

// TestTermGlyphsKeepEveryDistinctStrip checks a term whose nodes draw more than
// one picture keeps them all: the Archive nodes are all "Archive", and showing one
// of them would teach a player only half the vocabulary. At least one term carries
// several strips, and no term repeats one.
func TestTermGlyphsKeepEveryDistinctStrip(t *testing.T) {
	several := 0
	for term, strips := range termGlyphs() {
		seen := map[string]bool{}
		for _, gs := range strips {
			key := legendKey(gs)
			if seen[key] {
				t.Errorf("term %q lists the same strip twice", term)
			}
			seen[key] = true
		}
		if len(strips) > 1 {
			several++
		}
	}
	if several == 0 {
		t.Error("no term carries more than one strip; the legend is picking one per term")
	}
}

// TestLegendStripsCarryNoQuantity checks the legend drops a catalogued node's own
// numbers: the "1" on a sample GainAember is the sample's, not the term's, so the
// picture a player looks up carries no count. A comparison mark may still print a
// threshold, which is vocabulary rather than a count.
func TestLegendStripsCarryNoQuantity(t *testing.T) {
	for term, strips := range termGlyphs() {
		for _, gs := range strips {
			for _, g := range gs {
				if g.qty != 0 {
					t.Errorf("term %q draws %s with the sample quantity %d",
						term, g.asset, g.qty)
				}
			}
		}
	}
}

// TestRulebookDrawsTermGlyphs draws the page and checks a term's heading carries
// its strip: the legend must be on the page, not merely buildable.
func TestRulebookDrawsTermGlyphs(t *testing.T) {
	h := app.HTMLString(NewRulebook().Render())

	if !strings.Contains(h, `class="doc-legend"`) {
		t.Fatal("/rulebook draws no glyph legend beside any heading")
	}
	// Archive is the worked case: several catalogued nodes name it, so its heading
	// carries several strips.
	i := strings.Index(h, `id="archive"`)
	if i < 0 {
		t.Fatal("/rulebook has no Archive heading to hang a legend on")
	}
	heading := h[i:]
	if end := strings.Index(heading, "</h3>"); end >= 0 {
		heading = heading[:end]
	}
	if n := strings.Count(heading, `class="doc-legend-strip"`); n < 2 {
		t.Errorf("the Archive heading draws %d strips, want the several its nodes draw", n)
	}
}
