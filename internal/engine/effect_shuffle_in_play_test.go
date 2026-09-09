package engine

import "testing"

// TestShuffleFriendlyCardsInPlayIntoDeckText covers the two-sentence render.
func TestShuffleFriendlyCardsInPlayIntoDeckText(t *testing.T) {
	want := "shuffle each friendly card in play into your deck. " +
		"Draw a card for each card shuffled into your deck this way"
	if got := (ShuffleFriendlyCardsInPlayIntoDeck{}).Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
}

// TestShuffleFriendlyCardsInPlayIntoDeckEmpty checks that with nothing in play the
// effect shuffles nothing and draws nothing.
func TestShuffleFriendlyCardsInPlayIntoDeckEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.Register(NewCard("timequake", Brobnar, Tactic, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ShuffleFriendlyCardsInPlayIntoDeck{}.Resolve(ctx)
	if g.State.Deck[0].Count != 0 || g.State.Hand[0].Count != 0 {
		t.Errorf("empty board: deck=%d hand=%d, want 0/0",
			g.State.Deck[0].Count, g.State.Hand[0].Count)
	}
}

// TestShuffleFriendlyCardsInPlayIntoDeckResolve checks every friendly card in play
// — a creature, its upgrade, and an artifact — is shuffled into the deck and the
// controller draws one card per card shuffled this way. An enemy creature is left
// untouched.
func TestShuffleFriendlyCardsInPlayIntoDeckResolve(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.Register(NewCard("timequake", Brobnar, Tactic, Common), 0)

	creature := g.AddToBattleline(testCreature("creature", 4), 0)
	upgrade := g.Register(NewCard("boon", Brobnar, Upgrade, Common), 0)
	g.AttachUpgrade(creature, upgrade)
	artifact := g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 4), 1)

	// Give the controller cards to draw the three shuffled cards back into hand.
	for i := 0; i < 3; i++ {
		g.State.Deck[0].add(g.Register(testCreature("filler", 1), 0))
	}
	handBefore := g.State.Hand[0].Count

	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ShuffleFriendlyCardsInPlayIntoDeck{}.Resolve(ctx)

	for _, id := range []LocalID{creature, upgrade, artifact} {
		if !g.State.Deck[0].contains(id) && !g.State.Hand[0].contains(id) {
			t.Errorf("%s should be in the deck or drawn to hand after the shuffle", g.Name(id))
		}
	}
	if g.inPlay(creature) || g.inPlay(artifact) {
		t.Error("friendly creature and artifact should have left play")
	}
	if !g.inPlay(enemy) {
		t.Error("the enemy creature should have been left in play")
	}
	// Three cards were shuffled, so three are drawn; the deck of six loses three.
	if got := g.State.Hand[0].Count - handBefore; got != 3 {
		t.Errorf("drew %d cards, want 3", got)
	}
	if g.State.Deck[0].Count != 3 {
		t.Errorf("deck should hold 3 cards after the draw, count = %d", g.State.Deck[0].Count)
	}
}
