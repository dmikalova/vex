package engine

import "testing"

// TestEndTurnEndsTheTurn covers the EndTurn effect: it ends the active player's
// turn the moment it resolves, jumping to the end-of-turn phase without running
// the ready, draw, or end-of-turn steps, the way Book of leQ ends your turn on a
// Star Alliance reveal.
func TestEndTurnEndsTheTurn(t *testing.T) {
	if got := (EndTurn{}).Text(); got != "end your turn" {
		t.Errorf("Text = %q, want %q", got, "end your turn")
	}

	g := started(t)
	if g.State.Phase != PhasePlay {
		t.Fatalf("phase before EndTurn = %v, want PhasePlay", g.State.Phase)
	}
	EndTurn{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.Phase != PhaseEndOfTurn {
		t.Errorf("phase after EndTurn = %v, want PhaseEndOfTurn", g.State.Phase)
	}
}

// TestEndTurnNowIgnoresNonPlayPhase covers the guard: EndTurnNow only ends the
// turn from within the play phase, so a second call once the turn is already
// ending does nothing.
func TestEndTurnNowIgnoresNonPlayPhase(t *testing.T) {
	g := started(t)
	g.EndTurnNow()
	if g.State.Phase != PhaseEndOfTurn {
		t.Fatalf("phase after first EndTurnNow = %v, want PhaseEndOfTurn", g.State.Phase)
	}
	// A second call is a no-op: the phase is no longer PhasePlay.
	g.EndTurnNow()
	if g.State.Phase != PhaseEndOfTurn {
		t.Errorf("phase after second EndTurnNow = %v, want PhaseEndOfTurn", g.State.Phase)
	}
}

// TestEndTurnNowIgnoresWonGame covers the won-game guard: a decided game never
// ends another turn.
func TestEndTurnNowIgnoresWonGame(t *testing.T) {
	g := started(t)
	g.State.Winner = 0
	g.EndTurnNow()
	if g.State.Phase != PhasePlay {
		t.Errorf("phase after EndTurnNow on a won game = %v, want PhasePlay", g.State.Phase)
	}
}

// TestBookOfLeQComposition renders the composed Book of leQ ability, checking the
// reveal, active-house, and end-turn nodes fold into the printed rules text.
func TestBookOfLeQComposition(t *testing.T) {
	a := Ability{Trigger: TriggerAction, Effect: Sentences{Effects: []Effect{
		RevealTopOfDeck{Amount: 1},
		Conditional{
			Cond: ItIs{House: exceptHouse(StarAlliance)},
			Then: ChangeActiveHouse{To: TheContextualHouse},
			Else: EndTurn{},
		},
	}}}
	const want = "Action: Reveal the top card of your deck. If it is a non-Star Alliance card, its house becomes your active house. Otherwise, end your turn."
	if got := RenderAbility(a); got != want {
		t.Errorf("render =\n %q\nwant\n %q", got, want)
	}
}
