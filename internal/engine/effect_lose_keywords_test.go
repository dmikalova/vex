package engine

import "testing"

// TestLoseKeywordsText covers the printed clause for each supported duration and
// for one and two keywords.
func TestLoseKeywordsText(t *testing.T) {
	one := LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Elusive},
		Duration: RemainderOfPlayerTurn,
	}
	if got := one.Text(); got != "for the remainder of the turn, it loses elusive" {
		t.Errorf("text = %q", got)
	}
	two := LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Taunt, Elusive},
		Duration: RemainderOfPlayerTurn,
	}
	if got := two.Text(); got != "for the remainder of the turn, it loses taunt and elusive" {
		t.Errorf("text = %q", got)
	}
	self := LoseKeywords{
		Target:   Target{Kind: TargetThisCreature},
		Keywords: []Keyword{Elusive},
		Duration: StartOfPlayerNextTurn,
	}
	if got := self.Text(); got != "until the start of your next turn, "+SelfName+" loses elusive" {
		t.Errorf("next-turn text = %q", got)
	}
}

// TestLoseKeywordsValidate rejects a missing target, an empty keyword list, an
// unset keyword, and an unsupported duration.
func TestLoseKeywordsValidate(t *testing.T) {
	if err := (LoseKeywords{Keywords: []Keyword{Elusive}, Duration: RemainderOfPlayerTurn}).validate(); err == nil {
		t.Error("a missing target should be rejected")
	}
	if err := (LoseKeywords{Target: Target{Kind: TargetTriggeringCreature}, Duration: RemainderOfPlayerTurn}).validate(); err == nil {
		t.Error("an empty keyword list should be rejected")
	}
	if err := (LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{0},
		Duration: RemainderOfPlayerTurn,
	}).validate(); err == nil {
		t.Error("an unset keyword should be rejected")
	}
	if err := (LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Elusive},
	}).validate(); err == nil {
		t.Error("an unset duration should be rejected")
	}
	if err := (LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Elusive},
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("an unsupported duration should be rejected")
	}
	if err := (LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Taunt, Elusive},
		Duration: RemainderOfPlayerTurn,
	}).validate(); err != nil {
		t.Errorf("a valid strip should pass: %v", err)
	}
}

// TestLoseKeywordsResolveForTurn strips a creature of taunt and elusive for the
// turn, and the ready phase restores them.
func TestLoseKeywordsResolveForTurn(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("warded", 3, WithKeywords(Elusive, Taunt)), 0)
	if !g.hasKeyword(id, Elusive) || !g.hasKeyword(id, Taunt) {
		t.Fatal("the creature should start with elusive and taunt")
	}

	LoseKeywords{
		Target:   Target{Kind: TargetTriggeringCreature},
		Keywords: []Keyword{Taunt, Elusive},
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, It: id, HasIt: true})

	if g.hasKeyword(id, Elusive) {
		t.Error("elusive should be lost for the turn")
	}
	if g.hasKeyword(id, Taunt) {
		t.Error("taunt should be lost for the turn")
	}

	// A second loss of a keyword already lost is a no-op.
	before := len(g.Log)
	g.SetRecording(true)
	g.LoseKeywordFrom(id, Taunt)
	g.SetRecording(false)
	if len(g.Log) != before {
		t.Error("re-losing a keyword already lost should log nothing new")
	}

	g.EndPlayPhase(0)
	if !g.hasKeyword(id, Elusive) || !g.hasKeyword(id, Taunt) {
		t.Error("the lost keywords should return when the turn ends")
	}
}

// TestLoseKeywordsResolveNextTurn strips a creature of elusive until the
// controller's next turn: the loss survives the opponent's turn and lifts only at
// the start of the controller's own next turn. A second loss of the held keyword is
// a no-op.
func TestLoseKeywordsResolveNextTurn(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("rizzo", 1, WithKeywords(Elusive)), 0)
	if !g.hasKeyword(id, Elusive) {
		t.Fatal("the creature should start with elusive")
	}

	LoseKeywords{
		Target:   Target{Kind: TargetThisCreature},
		Keywords: []Keyword{Elusive},
		Duration: StartOfPlayerNextTurn,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: id})

	if g.hasKeyword(id, Elusive) {
		t.Error("elusive should be lost")
	}

	// A second loss of the held keyword logs nothing new.
	before := len(g.Log)
	g.SetRecording(true)
	g.LoseKeywordUntilNextTurn(id, Elusive)
	g.SetRecording(false)
	if len(g.Log) != before {
		t.Error("re-losing a keyword already lost should log nothing new")
	}

	// Neither ready phase lifts it, nor the opponent's start-of-turn; only the
	// controller's own start-of-turn restores the keyword.
	g.readyPhase(0)
	g.startOfTurnPhase(1)
	if g.hasKeyword(id, Elusive) {
		t.Error("the loss should survive the opponent's turn")
	}
	g.startOfTurnPhase(0)
	if !g.hasKeyword(id, Elusive) {
		t.Error("the loss should lift at the start of the controller's next turn")
	}
}
