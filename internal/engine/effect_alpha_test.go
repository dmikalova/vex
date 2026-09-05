package engine

import "testing"

// These tests cover the Alpha keyword: a card with Alpha must be the first card
// its player plays this turn.

func TestActedThisTurn(t *testing.T) {
	g := started(t)
	// A creature in play that has not been used does not count as having acted.
	idle := g.AddToBattleline(testCreature("idle", 3), 0)
	if g.actedThisTurn(0) {
		t.Fatal("a fresh turn with only an unused creature should not count as acted")
	}
	// Using that creature does count.
	g.State.Cards[idle].TimesUsedThisTurn = 1
	if !g.actedThisTurn(0) {
		t.Error("using a creature should count as having acted this turn")
	}
}

func TestAlphaMustBeFirst(t *testing.T) {
	g := started(t)
	alpha := g.AddToHand(
		NewCard("Ambush", Brobnar, Tactic, Common, WithKeywords(Alpha)),
		0,
	)
	// As the first thing done this turn, the Alpha card is playable.
	if err := g.CanPlay(0, alpha); err != nil {
		t.Fatalf("Alpha played first = %v, want nil", err)
	}

	// Once another card has been played, the Alpha card is barred.
	other := g.AddToHand(testCreature("c", 3), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, other), false); err != nil {
		t.Fatalf("play other card: %v", err)
	}
	if err := g.CanPlay(0, alpha); err != ErrAlphaNotFirst {
		t.Errorf("Alpha after acting (CanPlay) = %v, want ErrAlphaNotFirst", err)
	}
	if err := g.PlayAction(0, handIdxByID(g, 0, alpha)); err != ErrAlphaNotFirst {
		t.Errorf("Alpha after acting (PlayAction) = %v, want ErrAlphaNotFirst", err)
	}
}

func TestAlphaUpgradeBarred(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("host", 3), 0)
	up := g.AddToHand(
		NewCard("Vanguard", Brobnar, Upgrade, Common, WithKeywords(Alpha)),
		0,
	)
	first := g.AddToHand(testCreature("c", 3), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
		t.Fatalf("play first card: %v", err)
	}
	if _, err := g.PlayUpgrade(0, handIdxByID(g, 0, up)); err != ErrAlphaNotFirst {
		t.Errorf("Alpha upgrade after acting = %v, want ErrAlphaNotFirst", err)
	}
}
