package engine

import "testing"

// TestSelfSelectionText covers the Self selection's rendered fragments: it names
// the acting card by the {self} placeholder, with no article.
func TestSelfSelectionText(t *testing.T) {
	if got := (Self{}).noun(); got != SelfName {
		t.Errorf("Self.noun() = %q, want %q", got, SelfName)
	}
	if got := (Self{}).object(); got != SelfName {
		t.Errorf("Self.object() = %q, want %q", got, SelfName)
	}
}

// TestChosenAnotherExcludesTheCardInContext covers Resurgence's second pick: the
// word "another" is a promise, so an Another Chosen both renders it and drops
// ctx.It from the candidates rather than trusting the first pick to have moved
// the card out of the pile.
func TestChosenAnotherExcludesTheCardInContext(t *testing.T) {
	sel := Chosen{Type: Creature, Another: true}
	if got := sel.object(); got != "another creature" {
		t.Errorf("object() = %q, want %q", got, "another creature")
	}
	if sel.plainType() {
		t.Error("an Another Chosen is not a bare type and must not fold into a noun list")
	}

	g := NewGame("A", "B", 1)
	first := g.AddToDiscard(testCreature("first", 2), 0)
	second := g.AddToDiscard(testCreature("second", 2), 0)
	ctx := &EffectContext{Resolver: g, It: first, HasIt: true}

	got := sel.candidates(ctx, []LocalID{first, second})
	if len(got) != 1 || got[0] != second {
		t.Errorf("candidates = %v, want only the creature the first pick did not take", got)
	}
	if plain := (Chosen{Type: Creature}).candidates(
		ctx,
		[]LocalID{first, second},
	); len(
		plain,
	) != 2 {
		t.Errorf("without Another both creatures stay eligible, got %v", plain)
	}
}

// TestTriggersFromDiscardReturnsSelf covers the discard-trigger path end to end: a
// card with WithTriggersFromDiscard keeps its AfterChooseHouse ability live in its
// owner's discard pile, and a PutFromDiscard{Self} returns that very card to hand.
// This exercises the choose-house discard scan (game_turn.go), the ADR-0030 guard
// exception (resolveTriggered via activeInDiscard), and the Self selection's
// candidates/pick.
func TestTriggersFromDiscardReturnsSelf(t *testing.T) {
	creeper := NewCard("Creeper", Dis, Creature, Uncommon, WithPower(2),
		WithTriggersFromDiscard(),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Dis},
			Then: PutFromDiscard{Selection: Self{}, Destination: ToHand},
		}))

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	id := g.AddToDiscard(creeper, 0)

	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if !containsID(g.Hand(0), id) {
		t.Fatalf("Creeper should have returned to hand, hand=%v discard=%v",
			g.Hand(0), g.Discard(0))
	}
	if containsID(g.Discard(0), id) {
		t.Errorf("Creeper should have left the discard pile")
	}
}

// TestTriggersFromDiscardStaysWhenOtherHouseChosen covers the gate: choosing a
// different house does not fire the discard-pile ability, so the card stays put.
func TestTriggersFromDiscardStaysWhenOtherHouseChosen(t *testing.T) {
	creeper := NewCard("Creeper", Dis, Creature, Uncommon, WithPower(2),
		WithTriggersFromDiscard(),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Dis},
			Then: PutFromDiscard{Selection: Self{}, Destination: ToHand},
		}))

	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Dis, Logos, Untamed})
	g.StartTurn(0)
	id := g.AddToDiscard(creeper, 0)

	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if !containsID(g.Discard(0), id) {
		t.Errorf("Creeper should have stayed in the discard pile, discard=%v hand=%v",
			g.Discard(0), g.Hand(0))
	}
}

// TestSelfSelectionSkipsWhenSourceElsewhere covers Self.candidates and pick when
// the source card is not among the zone's cards: nothing is picked.
func TestSelfSelectionSkipsWhenSourceElsewhere(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	src := g.AddToDiscard(testCreature("src", 3), 0)
	other := g.AddToDiscard(testCreature("other", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if got := (Self{}).pick(ctx, []LocalID{other}); len(got) != 0 {
		t.Errorf("pick without the source present = %v, want none", got)
	}
	if got := (Self{}).candidates(ctx, []LocalID{src, other}); len(got) != 1 || got[0] != src {
		t.Errorf("candidates = %v, want [%v]", got, src)
	}
}
