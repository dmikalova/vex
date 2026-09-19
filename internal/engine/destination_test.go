package engine

import "testing"

// TestMoveFromDeckRowsAreDistinct pins the two deck-source rows of the move matrix
// that deck routing used to reach through its own table of resolver calls
// (DeckDest.mover). Both are easy to lose: a discard from a deck falls through to
// the hand arm if its case is dropped, and the bottom of the deck has no other
// source at all, so without its arm the card would silently land in a hand.
func TestMoveFromDeckRowsAreDistinct(t *testing.T) {
	for _, tc := range []struct {
		name string
		dest Destination
		want func(*Game, LocalID) bool
	}{
		{
			name: "discard",
			dest: toDiscard,
			want: func(g *Game, id LocalID) bool { return g.State.Discard[0].contains(id) },
		},
		{
			name: "bottom of deck",
			dest: ToBottomOfDeck,
			want: func(g *Game, id LocalID) bool {
				deck := g.State.Deck[0].slice()
				return len(deck) > 0 && deck[len(deck)-1] == id
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGame("A", "B", 1)
			top := g.Register(NewCard("top", Untamed, Creature, Common, WithPower(3)), 0)
			rest := g.Register(NewCard("rest", Untamed, Creature, Common, WithPower(3)), 0)
			g.State.Deck[0].add(top)
			g.State.Deck[0].add(rest)

			ctx := &EffectContext{Resolver: g, Controller: 0}
			tc.dest.moveFrom(ctx, Deck, 0, top)

			if g.State.Hand[0].contains(top) {
				t.Fatal("a deck-source move must not fall through to the hand arm")
			}
			if !tc.want(g, top) {
				t.Errorf("the card did not reach its destination")
			}
		})
	}
}
