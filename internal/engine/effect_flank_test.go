package engine

import "testing"

func TestConsiderFlankValidate(t *testing.T) {
	if err := (ConsiderFlank{}).validate(); err == nil {
		t.Error("an unset target should be rejected")
	}
	if err := (ConsiderFlank{Target: Target{Kind: TargetThisCreature}}).validate(); err != nil {
		t.Errorf("valid ConsiderFlank = %v", err)
	}
}

func TestConsiderFlankText(t *testing.T) {
	e := ConsiderFlank{Target: Target{Kind: TargetTriggeringCreature}}
	want := "for the remainder of the turn, it is considered a flank creature"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// A creature stuck in the middle of the battleline is not on a flank, but
// ConsiderFlank makes it count as one for every flank check.
func TestConsiderFlankResolve(t *testing.T) {
	g := started(t)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	mid := g.AddToBattleline(testCreature("mid", 3), 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	_, _ = left, right

	if g.onFlankOf(mid) {
		t.Fatal("mid creature should not start on a flank")
	}
	ctx := &EffectContext{Resolver: g, Source: mid, Controller: 0}
	if onFlank(ctx, mid) {
		t.Fatal("mid creature should not start on a flank (target filter)")
	}

	ConsiderFlank{Target: Target{Kind: TargetThisCreature}}.Resolve(ctx)

	if !g.ConsideredFlank(mid) {
		t.Error("mid creature should be considered a flank creature after resolve")
	}
	if !g.onFlankOf(mid) {
		t.Error("combat flank check should count the mid creature as a flank")
	}
	if !onFlank(ctx, mid) {
		t.Error("target flank filter should count the mid creature as a flank")
	}
}

// Considering an already-considered creature a flank again is a no-op — it neither
// re-records nor changes anything.
func TestConsiderFlankDedup(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("c", 3), 0)

	g.ConsiderFlank(id)
	before := len(g.Log)
	g.ConsiderFlank(id)
	if got := len(g.Log) - before; got != 0 {
		t.Errorf("re-considering logged %d entries, want 0", got)
	}
}

// The ready phase lifts the flank override so it lasts only the turn.
func TestConsiderFlankClearedInReadyPhase(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("c", 3), 0)
	g.ConsiderFlank(id)
	if !g.ConsideredFlank(id) {
		t.Fatal("setup: creature should be considered a flank")
	}

	g.readyPhase(0)

	if g.ConsideredFlank(id) {
		t.Error("ready phase should clear the flank override")
	}
}
