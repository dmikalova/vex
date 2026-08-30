package engine

import "testing"

func TestDiscardTopOfDeckAndRevealHandForAember(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	top := g.AddToDeck(NewCard("Mars Top", Mars, Action, Common), 1)
	next := g.AddToDeck(NewCard("Logos Next", Logos, Action, Common), 1)
	g.AddToHand(NewCard("Mars Creature", Mars, Creature, Common, WithPower(1)), 1)
	g.AddToHand(NewCard("Mars Action", Mars, Action, Common), 1)
	g.AddToHand(NewCard("Dis Action", Dis, Action, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := DiscardTopOfDeckAndRevealHandForAember{Player: Opponent, Gainer: Controller}
	if got := e.Text(); got != "discard the top card of your opponent's deck and reveal their hand. You gain 1 Æmber for each card of the discarded card's house revealed this way." {
		t.Errorf("text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v", err)
	}

	e.Resolve(ctx)
	if discard := g.Discard(1); len(discard) != 1 || discard[0] != top {
		t.Errorf("opponent discard = %v, want top card %d", discard, top)
	}
	if deck := g.Deck(1); len(deck) != 1 || deck[0] != next {
		t.Errorf("opponent deck = %v, want next card %d", deck, next)
	}
	if got := g.Aember(0); got != 2 {
		t.Errorf("controller Æmber = %d, want 2", got)
	}
	if ctx.Revealed != 3 {
		t.Errorf("revealed = %d, want whole hand count 3", ctx.Revealed)
	}
}

func TestDiscardTopOfDeckAndRevealHandForAemberReciprocalText(t *testing.T) {
	// The reciprocal half spells its own effect out rather than pointing back at
	// the preceding one, so each half reads on its own terms.
	e := DiscardTopOfDeckAndRevealHandForAember{Player: Controller, Gainer: Opponent}
	want := "discard the top card of your deck and reveal your hand. Your opponent gains 1 Æmber for each card of the discarded card's house revealed this way."
	if got := e.Text(); got != want {
		t.Errorf("reciprocal text = %q, want %q", got, want)
	}
}

func TestDiscardTopOfDeckAndRevealHandForAemberEmptyDeckSkips(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.AddToHand(NewCard("Mars Action", Mars, Action, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0, Revealed: 7}

	DiscardTopOfDeckAndRevealHandForAember{Player: Opponent, Gainer: Controller}.Resolve(ctx)

	if got := g.Aember(0); got != 0 {
		t.Errorf("controller Æmber = %d, want 0", got)
	}
	if ctx.Revealed != 0 {
		t.Errorf("revealed = %d, want 0", ctx.Revealed)
	}
}

func TestDiscardTopOfDeckAndRevealHandForAemberValidate(t *testing.T) {
	if err := (DiscardTopOfDeckAndRevealHandForAember{Gainer: Controller}).validate(); err == nil {
		t.Error("missing player should be rejected")
	}
	if err := (DiscardTopOfDeckAndRevealHandForAember{Player: Opponent}).validate(); err == nil {
		t.Error("missing gainer should be rejected")
	}
}
