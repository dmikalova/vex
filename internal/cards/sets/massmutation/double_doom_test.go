package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Double Doom
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Put an enemy creature into its owner's hand. Your opponent discards a random card from their hand.
func TestDoubleDoom(t *testing.T) {
	t.Run("returns an enemy creature to hand, then a random card is discarded", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(DoubleDoom)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Play(DoubleDoom)

		// the opponent starts with an empty hand, so after the bounce the returned
		// creature is the only random-discard candidate: it leaves play and lands in
		// the discard pile, exercising both halves.
		h.Expect(foe).At(ct.Discard)
	})
}
