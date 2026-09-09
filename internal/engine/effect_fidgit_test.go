package engine

import "testing"

func TestPlayDiscardedTacticFromOpponentText(t *testing.T) {
	want := "discard a random card from your opponent's archives or the top " +
		"card of their deck. If that card is a Tactic, play it as if it were yours"
	if got := (PlayDiscardedTacticFromOpponent{}).Text(); got != want {
		t.Errorf("text = %q", got)
	}
}

// TestPlayDiscardedTacticFromOpponentArchives covers Fidgit's reap taking a
// Tactic out of the opponent's archives: the card is discarded and then played as
// the reaping player's own — its Play: ability resolves under their control while
// its owner stays the opponent.
func TestPlayDiscardedTacticFromOpponentArchives(t *testing.T) {
	g := started(t)
	tactic := g.AddToArchives(NewCard("Snatched", Logos, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 3})), 1)

	// The default chooser answers ChooseOption with index 0 — the archives source.
	PlayDiscardedTacticFromOpponent{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Archives[1].contains(tactic) {
		t.Error("the tactic should have left the opponent's archives")
	}
	if got := g.State.Aember[0]; got != 3 {
		t.Errorf("player 0 Æmber = %d, want 3 from playing the tactic as yours", got)
	}
	if got := g.owner(tactic); got != 1 {
		t.Errorf("owner = %d, want 1 (unchanged)", got)
	}
}

// TestPlayDiscardedTacticFromOpponentDeckNonTactic covers the deck-top source with
// a non-Tactic card: it is discarded but not played.
func TestPlayDiscardedTacticFromOpponentDeckNonTactic(t *testing.T) {
	g := started(t)
	creature := g.AddToDeck(NewCard("Topper", Logos, Creature, Common, WithPower(3),
		WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 3})), 1)

	g.SetChooser(0, optionPicker{idx: 1}) // the top of the deck
	PlayDiscardedTacticFromOpponent{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.inPlay(creature) {
		t.Error("a non-Tactic should be discarded, not played")
	}
	if got := g.State.Aember[0]; got != 0 {
		t.Errorf("player 0 Æmber = %d, want 0 (nothing played)", got)
	}
	if d := g.Discard(1); len(d) != 1 || d[0] != creature {
		t.Errorf("owner discard = %v, want the discarded card", d)
	}
}

// TestPlayDiscardedTacticFromOpponentEmptySource covers an empty chosen source:
// the reap discards nothing and plays nothing.
func TestPlayDiscardedTacticFromOpponentEmptySource(t *testing.T) {
	g := started(t)

	// The default chooser picks the archives source, which is empty.
	PlayDiscardedTacticFromOpponent{}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.State.Aember[0]; got != 0 {
		t.Errorf("player 0 Æmber = %d, want 0 (nothing to discard)", got)
	}
}
