package engine

import "testing"

// TestShuffleFriendlyCardsIntoDeckText covers the single-sentence render; the draw
// is a separate effect on the card (Timequake), no longer part of the shuffle.
func TestShuffleFriendlyCardsIntoDeckText(t *testing.T) {
	want := "shuffle each friendly card in play into your deck"
	if got := (ShuffleFriendlyCardsIntoDeck{}).Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
}

// TestShuffleFriendlyCardsIntoDeckEmpty checks that with nothing in play the effect
// shuffles nothing and tallies nothing.
func TestShuffleFriendlyCardsIntoDeckEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.Register(NewCard("timequake", Brobnar, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ShuffleFriendlyCardsIntoDeck{}.Resolve(ctx)
	if ctx.Produced.Moved != [2]int{0, 0} {
		t.Errorf("empty board: Moved = %v, want [0 0]", ctx.Produced.Moved)
	}
	if g.State.Deck[0].Count != 0 {
		t.Errorf("empty board: deck = %d, want 0", g.State.Deck[0].Count)
	}
}

// TestShuffleFriendlyCardsIntoDeckResolve checks every friendly card in play — a
// creature, its upgrade, and an artifact — is shuffled into the controller's deck
// and tallied under the controller. An enemy creature is left untouched.
func TestShuffleFriendlyCardsIntoDeckResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.Register(NewCard("timequake", Brobnar, Tactic, Common), 0)

	creature := g.AddToBattleline(testCreature("creature", 4), 0)
	upgrade := g.Register(NewCard("boon", Brobnar, Upgrade, Common), 0)
	g.AttachUpgrade(creature, upgrade)
	artifact := g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 4), 1)

	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ShuffleFriendlyCardsIntoDeck{}.Resolve(ctx)

	for _, id := range []LocalID{creature, upgrade, artifact} {
		if !g.State.Deck[0].contains(id) {
			t.Errorf("%s should be in the controller's deck after the shuffle", g.Name(id))
		}
	}
	if g.inPlay(creature) || g.inPlay(artifact) {
		t.Error("friendly creature and artifact should have left play")
	}
	if !g.inPlay(enemy) {
		t.Error("the enemy creature should have been left in play")
	}
	if ctx.Produced.Moved != [2]int{3, 0} {
		t.Errorf("Moved = %v, want [3 0]", ctx.Produced.Moved)
	}
	if got := (CardsShuffledIntoDeck{}).Value(ctx); got != 3 {
		t.Errorf("CardsShuffledIntoDeck.Value = %d, want 3", got)
	}
	if got := (CardsShuffledIntoDeck{}).CountText(); got != "card shuffled into your deck this way" {
		t.Errorf("CountText = %q", got)
	}
}

// TestShuffleFriendlyCardsIntoDeckNotOwned checks that a card the controller plays
// but does not own is shuffled into its owner's deck and tallied under that owner,
// not the controller — so the controller draws only for the cards that returned to
// their own deck.
func TestShuffleFriendlyCardsIntoDeckNotOwned(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.Register(NewCard("timequake", Brobnar, Tactic, Common), 0)

	mine := g.AddToBattleline(testCreature("mine", 4), 0)
	// A creature the opponent owns but the controller controls (in their battleline).
	theirs := g.Register(testCreature("theirs", 4), 1)
	g.State.Battleline[0].add(theirs)

	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ShuffleFriendlyCardsIntoDeck{}.Resolve(ctx)

	if !g.State.Deck[0].contains(mine) {
		t.Error("the owned creature should return to the controller's deck")
	}
	if !g.State.Deck[1].contains(theirs) {
		t.Error("the enemy-owned creature should return to its owner's deck")
	}
	if ctx.Produced.Moved != [2]int{1, 1} {
		t.Errorf("Moved = %v, want [1 1]", ctx.Produced.Moved)
	}
	if got := (CardsShuffledIntoDeck{}).Value(ctx); got != 1 {
		t.Errorf("CardsShuffledIntoDeck.Value = %d, want 1 (only the owned card)", got)
	}
}
