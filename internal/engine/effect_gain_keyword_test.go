package engine

import "testing"

// TestGainKeywordsText covers the printed clause for each supported duration.
func TestGainKeywordsText(t *testing.T) {
	next := GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Keywords: []Keyword{Elusive},
		Duration: StartOfPlayerNextTurn,
	}
	if got := next.Text(); got != "each friendly creature gains elusive until the start of your next turn" {
		t.Errorf("next-turn text = %q", got)
	}
	turn := GainKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Skirmish},
		Duration: RemainderOfPlayerTurn,
	}
	if got := turn.Text(); got != "for the remainder of the turn, it gains skirmish" {
		t.Errorf("remainder text = %q", got)
	}
}

// TestGainKeywordsValidate rejects a missing target, an empty keyword list, an
// unset keyword, and an unsupported duration.
func TestGainKeywordsValidate(t *testing.T) {
	if err := (GainKeywords{
		Keywords: []Keyword{Elusive},
		Duration: StartOfPlayerNextTurn,
	}).validate(); err == nil {
		t.Error("a missing target should be rejected")
	}
	if err := (GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Duration: StartOfPlayerNextTurn,
	}).validate(); err == nil {
		t.Error("an empty keyword list should be rejected")
	}
	if err := (GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Keywords: []Keyword{0},
		Duration: StartOfPlayerNextTurn,
	}).validate(); err == nil {
		t.Error("an unset keyword should be rejected")
	}
	if err := (GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Keywords: []Keyword{Elusive},
	}).validate(); err == nil {
		t.Error("an unset duration should be rejected")
	}
	if err := (GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Keywords: []Keyword{Elusive},
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("an unsupported duration should be rejected")
	}
	if err := (GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Keywords: []Keyword{Elusive},
		Duration: StartOfPlayerNextTurn,
	}).validate(); err != nil {
		t.Errorf("a valid grant should pass: %v", err)
	}
}

// TestGainKeywordsResolveNextTurn grants each friendly creature the keyword until
// the start of the controller's next turn; it survives the opponent's turn.
func TestGainKeywordsResolveNextTurn(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	one := g.AddToBattleline(testCreature("one", 3), 0)
	two := g.AddToBattleline(testCreature("two", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	GainKeywords{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Keywords: []Keyword{Elusive},
		Duration: StartOfPlayerNextTurn,
	}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})

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

// TestGainKeywordsResolveForTurn grants the keyword for the remainder of the turn;
// the ready phase clears it (unlike the next-turn duration, which survives the
// opponent's turn).
func TestGainKeywordsResolveForTurn(t *testing.T) {
	g := started(t)
	one := g.AddToBattleline(testCreature("one", 3), 0)
	two := g.AddToBattleline(testCreature("two", 3), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         one,
		HasIt:      true,
	}

	GainKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Skirmish},
		Duration: RemainderOfPlayerTurn,
	}.Resolve(ctx)
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
