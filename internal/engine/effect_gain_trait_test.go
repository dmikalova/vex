package engine

import "testing"

// TestGainTraitText covers the printed clause, standalone.
func TestGainTraitText(t *testing.T) {
	e := GainTrait{Target: Target{Kind: TargetTriggeringCreature}, Trait: Mutant}
	want := "it gains the Mutant trait until the start of your next turn"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestGainTraitValidate rejects a missing target and an unset trait.
func TestGainTraitValidate(t *testing.T) {
	if err := (GainTrait{Trait: Mutant}).validate(); err == nil {
		t.Error("a missing target should be rejected")
	}
	if err := (GainTrait{Target: Target{Kind: TargetTriggeringCreature}}).validate(); err == nil {
		t.Error("an unset trait should be rejected")
	}
	if err := (GainTrait{
		Target: Target{Kind: TargetTriggeringCreature},
		Trait:  Mutant,
	}).validate(); err != nil {
		t.Errorf("a valid grant should pass: %v", err)
	}
}

// TestGainTraitResolve grants a creature the trait until the controller's next
// turn: it survives the opponent's turn and lifts only at the start of the
// controller's own next turn. A second grant of the held trait is a no-op, and granting to a
// creature no longer in play does not panic.
func TestGainTraitResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	beast := g.AddToBattleline(testCreature("beast", 4), 0)

	GainTrait{Target: Target{Kind: TargetTriggeringCreature}, Trait: Mutant}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, It: beast, HasIt: true})

	if !g.HasTrait(beast, Mutant) {
		t.Error("the creature should have gained the Mutant trait")
	}
	if !hasLogLine(g, "beast gains the Mutant trait") {
		t.Error("the grant should be logged")
	}
	// Asking for no trait is never a match, even on a creature that has a granted
	// trait.
	if g.HasTrait(beast, traitUnset) {
		t.Error("an unset trait should never match")
	}

	// A second grant of the held trait logs nothing new.
	before := len(g.Log)
	g.GrantTraitUntilNextTurn(beast, Mutant)
	if len(g.Log) != before {
		t.Error("re-granting a held trait should log nothing new")
	}

	// Neither ready phase lifts it, and the opponent's start-of-turn does not; only
	// the controller's own start-of-turn does, before start-of-turn abilities run.
	g.readyPhase(0)
	g.startOfTurnPhase(1)
	if !g.HasTrait(beast, Mutant) {
		t.Error("the grant should survive the opponent's turn")
	}
	g.startOfTurnPhase(0)
	if g.HasTrait(beast, Mutant) {
		t.Error("the grant should lift at the start of the controller's next turn")
	}

	// Granting to a creature no longer in play is a safe no-op.
	gone := g.AddToBattleline(testCreature("gone", 3), 0)
	g.DestroyEach(0, []LocalID{gone})
	g.GrantTraitUntilNextTurn(gone, Mutant)
}
