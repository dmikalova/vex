package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestAemberEffectGlyphsTranscription binds the æmber family to a worked example
// so the grammar cannot drift silently.
func TestAemberEffectGlyphsTranscription(t *testing.T) {
	// A Per variant of GainAember maps to the Æmber glyph without a numeral: the
	// count scales by a board quantity that lives in the text, not on the glyph.
	gs, covered := effectGlyphs(engine.GainAember{
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
