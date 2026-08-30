package engine

import "testing"

func TestDiscardTopOfDeckPlayer(t *testing.T) {
	// Text renders from each perspective: the granted default names "its
	// controller", while Controller and Opponent are direct first/second person.
	if got := (DiscardTopOfDeck{}).Text(); got != "discard the top card of its controller's deck" {
		t.Errorf("granted text = %q", got)
	}
	if got := (DiscardTopOfDeck{Player: Controller}).Text(); got != "discard the top card of your deck" {
		t.Errorf("controller text = %q", got)
	}
	if got := (DiscardTopOfDeck{Player: Opponent}).Text(); got != "discard the top card of your opponent's deck" {
		t.Errorf("opponent text = %q", got)
	}

	g := NewGame("Alice", "Bob", 1)
	top := g.AddToDeck(NewCard("Opp Top", Mars, Tactic, Common), 1)
	next := g.AddToDeck(NewCard("Opp Next", Logos, Tactic, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	DiscardTopOfDeck{Player: Opponent}.Resolve(ctx)
	if !ctx.HasIt || ctx.It != top {
		t.Errorf("context card = %d (has %v), want opponent top %d", ctx.It, ctx.HasIt, top)
	}
	if discard := g.Discard(1); len(discard) != 1 || discard[0] != top {
		t.Errorf("opponent discard = %v, want top %d", discard, top)
	}
	if deck := g.Deck(1); len(deck) != 1 || deck[0] != next {
		t.Errorf("opponent deck = %v, want next %d", deck, next)
	}
}

func TestDiscardTopOfDeckEmpty(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0, HasIt: true, It: 42}
	DiscardTopOfDeck{Player: Controller}.Resolve(ctx)
	if ctx.HasIt {
		t.Errorf("empty deck should leave no card in context, got It=%d", ctx.It)
	}
}

func TestCardsRevealedOfItsHouse(t *testing.T) {
	if got := (CardsRevealedOfItsHouse{}).CountText(); got != "card of the discarded card's house revealed this way" {
		t.Errorf("count text = %q", got)
	}

	g := NewGame("Alice", "Bob", 1)
	discarded := g.AddToDeck(NewCard("Mars Top", Mars, Tactic, Common), 1)
	g.AddToHand(NewCard("Mars Creature", Mars, Creature, Common, WithPower(1)), 1)
	g.AddToHand(NewCard("Mars Action", Mars, Tactic, Common), 1)
	g.AddToHand(NewCard("Dis Action", Dis, Tactic, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// With no card in context an empty deck established no house — count zero.
	if got := (CardsRevealedOfItsHouse{Player: Opponent}).Value(ctx); got != 0 {
		t.Errorf("no-context value = %d, want 0", got)
	}

	// Discard sets the context card, so the count matches the Mars hand cards.
	DiscardTopOfDeck{Player: Opponent}.Resolve(ctx)
	if ctx.It != discarded {
		t.Fatalf("context card = %d, want %d", ctx.It, discarded)
	}
	if got := (CardsRevealedOfItsHouse{Player: Opponent}).Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2 Mars cards", got)
	}
}
