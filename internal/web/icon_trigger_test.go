package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the trigger
// glyph vocabulary in icon_trigger.go.

// TestTriggerIconCoversEveryTrigger walks engine.Triggers() and fails on a
// trigger that falls through to the abstract unknown glyph. The card walk only
// reaches the triggers an implemented card happens to print, so a trigger no card
// uses yet would ship undrawn; the census guards Triggers() complete against the
// source, which makes this the trigger head's totality gate (ADR 0022).
func TestTriggerIconCoversEveryTrigger(t *testing.T) {
	for _, tr := range engine.Triggers() {
		name := triggerIcon(tr)
		if name == "" || name == "glyph-unknown" {
			t.Errorf("trigger %v draws %q; give it a glyph (ADR 0022)", tr, name)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("trigger %v: no web/assets/%s.svg", tr, name)
		}
	}
}

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
