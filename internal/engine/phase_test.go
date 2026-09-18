package engine

import "testing"

// TestPhaseTableIsTotal is the guard the phase descriptor table exists for: it
// proves the table describes exactly the phases Phases() lists, and that every
// entry is filled in. A phase added to the enum and to Phases() but not to the
// table fails here instead of silently rendering as "no phase", losing its
// rulebook heading, and doing nothing when it runs.
func TestPhaseTableIsTotal(t *testing.T) {
	for _, p := range Phases() {
		info, ok := phases[p]
		if !ok {
			t.Errorf("phase %d has no table entry", p)
			continue
		}
		if info.name == "" {
			t.Errorf("phase %d has no player-facing name", p)
		}
		if info.rulebookStep == "" {
			t.Errorf("phase %q has no rulebook heading", info.name)
		}
		// Every phase either blocks for a frontend decision or has a body the engine
		// runs; a phase that does neither would silently pass a turn by.
		if info.waitsForInput == (info.run != nil) {
			t.Errorf("phase %q must either wait for input or have a run body, not both or neither",
				info.name)
		}
	}
	if len(phases) != len(Phases()) {
		t.Errorf("table describes %d phases, Phases() lists %d", len(phases), len(Phases()))
	}
}

// The unset zero value is not a phase, so it renders as such, has no rulebook
// heading, waits for nothing, and runs nothing.
func TestUnsetPhaseIsNotAPhase(t *testing.T) {
	if phaseUnset.valid() {
		t.Error("the zero value should not be a valid phase")
	}
	if got := phaseUnset.String(); got != "no phase" {
		t.Errorf("String = %q, want %q", got, "no phase")
	}
	if got := phaseUnset.rulebookStep(); got != "" {
		t.Errorf("rulebookStep = %q, want empty", got)
	}
	if phaseUnset.waitsForInput() {
		t.Error("the unset phase should not wait for input")
	}
	g := started(t)
	g.State.Phase = phaseUnset
	g.runPhase()
}
