package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestTimedAssaultDrawsItsWindow pins the second pair the duration channel
// separates: engine.GainAssault grants the keyword for good and
// engine.GainAssaultUntilNextTurn lends it until the lender's next turn, and the
// two drew byte-identical strips until the turn marker was appended (ADR 0047).
//
// Declaring them synonyms was the alternative, on the grounds that a lasting buff
// and a permanent one are one act with a different shelf life. They are drawn
// apart instead: the marker costs one slot on one rare effect, and how long a buff
// lasts is a thing a player plays around.
func TestTimedAssaultDrawsItsWindow(t *testing.T) {
	target := engine.Target{Kind: engine.TargetThisCreature}
	timed := effectGlyphs(engine.GainAssaultUntilNextTurn{
		Target: target,
		Amount: engine.Fixed(2),
	})
	if len(timed) == 0 || timed[len(timed)-1].asset != "phase-turn" {
		t.Fatalf("glyphs = %+v, want the turn marker closing the run", timed)
	}
	kept := effectGlyphs(engine.GainAssault{Target: target, Amount: engine.Fixed(2)})
	if fmt.Sprintf("%+v", kept) == fmt.Sprintf("%+v", timed) {
		t.Errorf("the kept and timed assault grants still draw the same strip %+v", kept)
	}
}
