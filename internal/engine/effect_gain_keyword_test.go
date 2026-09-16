package engine

import "testing"

// TestGainKeywordText covers the printed clause.
func TestGainKeywordText(t *testing.T) {
	e := GainKeyword{Target: Target{Kind: TargetEachFriendlyCreature}, Keyword: Elusive}
	if got := e.Text(); got != "each friendly creature gains elusive until the start of your next turn" {
		t.Errorf("text = %q", got)
	}
}

// TestGainKeywordValidate rejects a missing target and an unset keyword.
func TestGainKeywordValidate(t *testing.T) {
	if err := (GainKeyword{Keyword: Elusive}).validate(); err == nil {
		t.Error("a missing target should be rejected")
	}
	if err := (GainKeyword{Target: Target{Kind: TargetEachFriendlyCreature}}).validate(); err == nil {
		t.Error("an unset keyword should be rejected")
	}
	if err := (GainKeyword{
		Target:  Target{Kind: TargetEachFriendlyCreature},
		Keyword: Elusive,
	}).validate(); err != nil {
		t.Errorf("a valid grant should pass: %v", err)
	}
}

// TestGainKeywordResolve grants each friendly creature the keyword; it survives
// the opponent's turn and lifts only at the start of the controller's own next turn.
func TestGainKeywordResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	one := g.AddToBattleline(testCreature("one", 3), 0)
	two := g.AddToBattleline(testCreature("two", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	GainKeyword{Target: Target{Kind: TargetEachFriendlyCreature}, Keyword: Elusive}.
		Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.hasKeyword(one, Elusive) || !g.hasKeyword(two, Elusive) {
		t.Error("both friendly creatures should have gained elusive")
	}
	if g.hasKeyword(enemy, Elusive) {
		t.Error("the enemy creature should not have gained elusive")
	}

	// A second grant of a keyword already held is a no-op that logs nothing.
	before := len(g.Log)
	g.GrantKeywordUntilNextTurn(one, Elusive)
	if len(g.Log) != before {
		t.Error("re-granting a held keyword should log nothing new")
	}

	// Neither ready phase lifts it, nor the opponent's start-of-turn.
	g.readyPhase(0)
	g.startOfTurnPhase(1)
	if !g.hasKeyword(one, Elusive) {
		t.Error("the grant should survive the opponent's turn")
	}

	// The controller's own next start-of-turn lifts it.
	g.startOfTurnPhase(0)
	if g.hasKeyword(one, Elusive) {
		t.Error("the grant should lift at the start of the controller's next turn")
	}
}
