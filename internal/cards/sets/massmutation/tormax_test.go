package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tormax
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Demon
//
//	Play/Fight/Reap: Discard your hand, and purge 2 random cards from your opponent's hand.
func TestTormax(t *testing.T) {
	var tormax, mine1, mine2, theirs1, theirs2 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand: ct.Cards(
				ct.Bind(&tormax, Tormax), card.GiganticArt(Tormax),
				ct.Bind(&mine1, ct.Tactic()),
				ct.Bind(&mine2, ct.Tactic()),
			),
		},
		P2: ct.Side{Hand: ct.Cards(
			ct.Bind(&theirs1, ct.Tactic()),
			ct.Bind(&theirs2, ct.Tactic()),
		)},
	})

	h.P1.Play(tormax)

	// Own hand discarded; the opponent's two cards purged.
	h.Expect(mine1).At(ct.Discard)
	h.Expect(mine2).At(ct.Discard)
	h.Expect(theirs1).At(ct.Purge)
	h.Expect(theirs2).At(ct.Purge)
}
