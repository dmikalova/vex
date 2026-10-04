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
	gs := effectGlyphs(engine.GainAember{
		Player: engine.Controller,
		Amount: 1,
		Per:    engine.ForgedKeys{Player: engine.Opponent},
	})
	if len(gs) != 1 || gs[0].asset != "aember" || gs[0].qty != 0 {
		t.Errorf(
			"GainAember Per glyphs = %+v, want a single æmber glyph with no numeral",
			gs,
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
			gs := effectGlyphs(tc.e)
			if len(gs) == 0 {
				t.Fatalf("%s drew no glyphs", tc.name)
			}
			got := gs[0]
			if got.asset != tc.asset || got.sigil != tc.sigil || got.qty != tc.qty {
				t.Errorf(
					"%s head glyph = %+v, want asset %q sigil %v qty %d",
					tc.name, got, tc.asset, tc.sigil, tc.qty,
				)
			}
		})
	}
}

// TestEveryLossDraws holds the Loss family's gate from the other side: every share
// of a pool an economy verb can move draws a real, distinct strip of existing
// assets, so "lose all your Æmber" cannot draw what "lose half your Æmber" draws.
func TestEveryLossDraws(t *testing.T) {
	assertFamilyDraws(t, "Loss")
	assertFamilyRowsDrawApart(t, "Loss")
}

// TestLossReachesTheVerbCarryingIt checks the share is drawn beside the verb that
// moves it: Too Much to Protect steals all but six, and before the share was drawn
// its strip was the plain steal mark — the same picture as stealing one Æmber.
func TestLossReachesTheVerbCarryingIt(t *testing.T) {
	gs := effectGlyphs(engine.StealAember{By: engine.AllBut{Keep: 6}})
	if len(gs) != 4 {
		t.Fatalf("share steal = %+v, want the steal and the share it leaves", gs)
	}
	if gs[2].qty != 6 {
		t.Errorf("share steal = %+v, want the pool of 6 it leaves behind", gs)
	}
	if gs[3].asset != "glyph-ban" {
		t.Errorf("share steal = %+v, want the refusal mark closing the remainder", gs)
	}
	plain := effectGlyphs(engine.StealAember{Amount: 1})
	if len(plain) != 1 {
		t.Errorf("a plain steal = %+v, want no share drawn at all", plain)
	}
}

// TestFractionsDrawTheirDenominator pins how a fractional share is numbered: a
// half draws a 2 and a third a 3, which is the only thing separating the two in a
// strip that speaks no words. Which way the share rounds is fine print and stays
// in the rules text, so the two halves draw alike on purpose.
func TestFractionsDrawTheirDenominator(t *testing.T) {
	for _, tc := range []struct {
		name string
		f    engine.Fraction
		want int
	}{
		{"half rounded down", engine.HalfRoundedDown, 2},
		{"half rounded up", engine.HalfRoundedUp, 2},
		{"third rounded down", engine.ThirdRoundedDown, 3},
		{"third rounded up", engine.ThirdRoundedUp, 3},
	} {
		gs := lossGlyphs(tc.f)
		if len(gs) != 1 {
			t.Fatalf("%s drew %+v, want one glyph", tc.name, gs)
		}
		if gs[0].qty != tc.want {
			t.Errorf("%s drew qty %d, want %d", tc.name, gs[0].qty, tc.want)
		}
	}
}
