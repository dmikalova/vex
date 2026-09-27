package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the trigger
// glyph vocabulary in icon_trigger.go.

// TestTriggerIconPhase checks the phase triggers (start/end of turn, end of ready
// step) share the turn-phase glyph rather than falling back to the unknown glyph.
func TestTriggerIconPhase(t *testing.T) {
	for _, tr := range []engine.Trigger{
		engine.TriggerStartOfTurn, engine.TriggerEndOfTurn, engine.TriggerEndOfReadyStep,
	} {
		if got := triggerIcon(tr); got != "phase-turn" {
			t.Errorf("triggerIcon(%v) = %q, want phase-turn", tr, got)
		}
	}
}
