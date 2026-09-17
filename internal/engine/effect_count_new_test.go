package engine

import "testing"

func TestArtifactsInPlayCount(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("a1", Brobnar, Artifact, Common), 0)
	g.AddArtifact(NewCard("a2", Brobnar, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := (ArtifactsInPlay{}).Value(ctx); got != 2 {
		t.Errorf("ArtifactsInPlay = %d, want 2", got)
	}
	if got := (ArtifactsInPlay{}).CountText(); got != "artifact in play" {
		t.Errorf("CountText = %q", got)
	}
}

func TestPowerCountersOnThisCount(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.AddPowerCounter(c, 4)
	ctx := &EffectContext{Resolver: g, Source: c, Controller: 0}

	if got := (PowerCountersOnThis{}).Value(ctx); got != 4 {
		t.Errorf("PowerCountersOnThis = %d, want 4", got)
	}
	want := "the number of +1 power counters on " + SelfName
	if got := (PowerCountersOnThis{}).CountText(); got != want {
		t.Errorf("CountText = %q, want %q", got, want)
	}
}

// TestCannotBeDealtDamageOpponentNextTurn covers the cross-turn immunity (Lucky
// Dice): dormant on the caster's turn, active during the opponent's next turn.
func TestCannotBeDealtDamageOpponentNextTurn(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 5), 0)
	e := CannotBeDealtDamage{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Duration: OpponentNextTurn,
	}
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.DamageImmune(c) {
		t.Error("OpponentNextTurn immunity should be dormant on the caster's turn")
	}
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if !g.DamageImmune(c) {
		t.Error("immunity should be active during the opponent's next turn")
	}
	want := "during your opponent's next turn, each friendly creature cannot be dealt damage"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}
