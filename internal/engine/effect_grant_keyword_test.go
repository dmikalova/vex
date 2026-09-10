package engine

import "testing"

// TestGainKeywordForTurnText covers the printed clause.
func TestGainKeywordForTurnText(t *testing.T) {
	e := GainKeywordForTurn{Target: Target{Kind: TargetTriggeringCreature}, Keyword: Skirmish}
	if got := e.Text(); got != "for the remainder of the turn, it gains skirmish" {
		t.Errorf("text = %q", got)
	}
}

// TestGainKeywordForTurnValidate rejects a missing target and an unset keyword.
func TestGainKeywordForTurnValidate(t *testing.T) {
	if (GainKeywordForTurn{Keyword: Skirmish}).validate() == nil {
		t.Error("unset target should be invalid")
	}
	if (GainKeywordForTurn{Target: Target{Kind: TargetTriggeringCreature}}).validate() == nil {
		t.Error("unset keyword should be invalid")
	}
	if (GainKeywordForTurn{
		Target:  Target{Kind: TargetTriggeringCreature},
		Keyword: Skirmish,
	}).validate() != nil {
		t.Error("a set target and keyword should be valid")
	}
}

// TestGainKeywordForTurnResolve grants the keyword for the turn; the ready phase
// clears it (unlike GainKeyword, which survives the opponent's turn).
func TestGainKeywordForTurnResolve(t *testing.T) {
	g := started(t)
	one := g.AddToBattleline(testCreature("one", 3), 0)
	two := g.AddToBattleline(testCreature("two", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: one, HasIt: true}

	GainKeywordForTurn{
		Target:  Target{Kind: TargetTriggeringCreature},
		Keyword: Skirmish,
	}.Resolve(
		ctx,
	)
	if !g.hasKeyword(one, Skirmish) {
		t.Error("the chosen creature should have gained skirmish")
	}
	if g.hasKeyword(two, Skirmish) {
		t.Error("only the chosen creature should have gained skirmish")
	}

	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.hasKeyword(one, Skirmish) {
		t.Error("the ready phase should clear the keyword gained for the turn")
	}
}
