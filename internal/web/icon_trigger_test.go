package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// This file is part of the Iconography pass (ADR 0022): tests for the trigger
// glyph vocabulary in icon_trigger.go.

// TestTriggerIconCoversEveryTrigger walks engine.Triggers() and fails on a
// trigger that falls through the switch and draws nothing. The card walk only
// reaches the triggers an implemented card happens to print, so a trigger no card
// uses yet would ship undrawn; the census guards Triggers() complete against the
// source, which makes this the trigger head's totality gate (ADR 0022).
//
// A trigger reaching the default arm panics here rather than returning "":
// triggerIcon reports an unmapped trigger through unmappedAsset, which the
// package's tests arm (icon_test.go). The empty check is what catches a trigger
// mapped to a blank stem by hand.
func TestTriggerIconCoversEveryTrigger(t *testing.T) {
	for _, tr := range engine.Triggers() {
		name := triggerIcon(tr)
		if name == "" {
			t.Errorf("trigger %v draws nothing; give it a glyph (ADR 0022)", tr)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("trigger %v: no web/assets/%s.svg", tr, name)
		}
	}
}

// TestTriggerIconPhase checks the phase triggers (start/end of turn, end of ready
// step) share the turn-phase glyph rather than falling through the switch.
func TestTriggerIconPhase(t *testing.T) {
	for _, tr := range []engine.Trigger{
		engine.TriggerStartOfTurn, engine.TriggerEndOfTurn, engine.TriggerEndOfReadyStep,
	} {
		if got := triggerIcon(tr); got != "phase-turn" {
			t.Errorf("triggerIcon(%v) = %q, want phase-turn", tr, got)
		}
	}
}
