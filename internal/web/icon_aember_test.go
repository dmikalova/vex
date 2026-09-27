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

// TestAemberVerbsTakeTheirOwnSigil pins the verb axis on the family it was worked
// out against (ADR 0047): gain, lose, steal, capture and exalt all sit on the same
// Æmber noun, so each carries its own verb mark instead of collapsing onto the
// others. It also holds the quantity slot free for the amount — LoseAember used to
// fake a modifier by writing "−" into it, the only effect in the pass that did.
func TestAemberVerbsTakeTheirOwnSigil(t *testing.T) {
	for _, tc := range []struct {
		name  string
		e     engine.Effect
		asset string
		sigil sigil
		qty   int
	}{
		{"gain", engine.GainAember{Player: engine.Controller, Amount: 2}, "aember", sigilGain, 2},
		{"lose", engine.LoseAember{Player: engine.Opponent, Amount: 3}, "aember", sigilLose, 3},
		{"steal", engine.StealAember{Amount: 1}, "aember", sigilSteal, 1},
		{"capture", engine.CaptureAember{Amount: 1}, "capture", sigilNone, 1},
		{"exalt", engine.Exalt{Amount: 1}, "aember", sigilExalt, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gs, covered := effectGlyphs(tc.e)
			if !covered || len(gs) == 0 {
				t.Fatalf("%s glyphs = %+v (covered=%v)", tc.name, gs, covered)
			}
			got := gs[0]
			if got.asset != tc.asset || got.sigil != tc.sigil || got.qty != tc.qty {
				t.Errorf(
					"%s head glyph = %+v, want asset %q sigil %v qty %d",
					tc.name, got, tc.asset, tc.sigil, tc.qty,
				)
			}
			if got.text != "" {
				t.Errorf(
					"%s writes %q into the text slot; the verb is the sigil's job",
					tc.name,
					got.text,
				)
			}
		})
	}
}
