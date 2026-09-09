package engine

import "testing"

// Quixxle Stone: while it is in play, any player who controls more creatures than
// their opponent cannot play creatures — symmetric, evaluated per attempting
// player, and only barring creatures.
func TestConditionalPlayBarText(t *testing.T) {
	def := NewCard(
		"Quixxle Stone",
		StarAlliance,
		Artifact,
		Rare,
		WithCannotPlayWhile(ConditionalPlayBar{Type: Creature, When: ControlsMoreCreatures{}}),
	)
	rules := cardRules(&def, false)
	want := "If a player has more creatures in play than their opponent, they cannot play creatures."
	found := false
	for _, r := range rules {
		if r == want {
			found = true
		}
	}
	if !found {
		t.Errorf("cardRules = %v, want a line %q", rules, want)
	}
}

// symmetricCondText falls back to the condition's own CondText (minus its "if "
// prefix) for any condition other than ControlsMoreCreatures.
func TestSymmetricCondTextFallback(t *testing.T) {
	if got := symmetricCondText(
		OpponentHasMoreKeys{},
	); got != "your opponent has more forged keys than you" {
		t.Errorf("symmetricCondText fallback = %q", got)
	}
}

func TestConditionalPlayBarBarsAheadPlayer(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	g.AddArtifact(
		NewCard(
			"Quixxle Stone",
			StarAlliance,
			Artifact,
			Rare,
			WithCannotPlayWhile(ConditionalPlayBar{Type: Creature, When: ControlsMoreCreatures{}}),
		),
		1, // controlled by the opponent, yet it bars whichever side is ahead
	)

	// Even counts (0 vs 0): the active player is not ahead, so creatures are legal.
	if g.cannotPlayCreatures(0) {
		t.Fatal("with equal creature counts, plays should not be barred")
	}

	// Player 0 pulls ahead: one creature to none.
	g.AddToBattleline(testCreature("ahead", 3), 0)
	if !g.cannotPlayCreatures(0) {
		t.Fatal("the player with more creatures should be barred from playing creatures")
	}
	g.AddToHand(testCreature("newbie", 2), 0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "newbie"), false); err != ErrCannotPlayCreature {
		t.Errorf("PlayCreature = %v, want ErrCannotPlayCreature", err)
	}

	// The bar only blocks creatures: a non-creature is still playable.
	g.AddToHand(NewCard("act", Brobnar, Tactic, Common), 0)
	if err := g.PlayAction(0, handIdx(g, 0, "act")); err != nil {
		t.Errorf("actions should still be playable: %v", err)
	}

	// The opponent draws even (1 vs 1): the active player is no longer ahead.
	g.AddToBattleline(testCreature("evener", 1), 1)
	if g.cannotPlayCreatures(0) {
		t.Fatal("with counts even again, the bar should lift")
	}
	g.AddToHand(testCreature("second", 2), 0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "second"), false); err != nil {
		t.Errorf("PlayCreature once even = %v, want nil", err)
	}
}
