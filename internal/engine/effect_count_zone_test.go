package engine

import "testing"

// TestCardsInZone covers the shared per-zone card count: its value in each zone,
// the possessive text for both players, the CountIs clause it renders, and the
// zone validation.
func TestCardsInZone(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.AddToDeck(NewCard("Deck 1", Mars, Tactic, Common), 0)
	g.AddToDeck(NewCard("Deck 2", Mars, Tactic, Common), 0)
	g.AddToHand(NewCard("Hand 1", Mars, Tactic, Common), 0)
	g.AddToArchives(NewCard("Arch 1", Mars, Tactic, Common), 0)
	g.AddToDiscard(NewCard("Disc 1", Mars, Tactic, Common), 0)
	g.AddToDiscard(NewCard("Disc 2", Mars, Tactic, Common), 0)
	g.AddToDiscard(NewCard("Disc 3", Mars, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	values := map[Zone]int{Deck: 2, Hand: 1, Archives: 1, Discard: 3}
	for zone, want := range values {
		if got := (CardsInZone{Zone: zone, Player: Controller}).Value(ctx); got != want {
			t.Errorf("Value(%v) = %d, want %d", zone.noun(), got, want)
		}
	}

	// CountText names the zone and whose it is.
	if got := (CardsInZone{Zone: Archives, Player: Controller}).CountText(); got != "card in your archives" {
		t.Errorf("friendly CountText = %q", got)
	}
	if got := (CardsInZone{Zone: Archives, Player: Opponent}).CountText(); got != "card in your opponent's archives" {
		t.Errorf("enemy CountText = %q", got)
	}

	// CountClause renders the "if" clause for both players.
	friendly := CardsInZone{Zone: Deck, Player: Controller}
	if got := friendly.CountClause(
		"5 or fewer",
		true,
	); got != "you have 5 or fewer cards in your deck" {
		t.Errorf("friendly CountClause = %q", got)
	}
	enemy := CardsInZone{Zone: Deck, Player: Opponent}
	if got := enemy.CountClause(
		"5 or more",
		true,
	); got != "your opponent has 5 or more cards in your opponent's deck" {
		t.Errorf("enemy CountClause = %q", got)
	}

	// validate accepts every nameable zone and rejects an unnamed one.
	for zone := range values {
		if err := (CardsInZone{Zone: zone}).validate(); err != nil {
			t.Errorf("valid zone %v rejected: %v", zone.noun(), err)
		}
	}
	if err := (CardsInZone{}).validate(); err == nil {
		t.Error("an unnamed zone should be rejected")
	}

	// A valid CountIs over a CardsInZone passes validation (AtMost is now allowed).
	if err := (CountIs{Count: friendly, Is: AtMost, Amount: 5}).validate(); err != nil {
		t.Errorf("valid AtMost CountIs rejected: %v", err)
	}

	// A Not may only wrap an AtLeast CountIs: AtMost has no "fewer than" negation.
	if err := (Not{Cond: CountIs{Count: friendly, Is: AtMost, Amount: 5}}).validate(); err == nil {
		t.Error("Not over a non-AtLeast CountIs should be rejected")
	}
}
