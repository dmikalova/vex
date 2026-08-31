package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Troop Call
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put each Niffle trait creature from your discard pile into your hand. Put each friendly Niffle trait creature into its owner's hand.
func TestTroopCall(t *testing.T) {
	t.Run("returns friendly Niffle creatures from discard and play to hand", func(t *testing.T) {
		var buried, inPlay, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Untamed,
				Hand:    ct.Cards(TroopCall),
				Discard: ct.Cards(ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits("Niffle")))),
				InPlay: ct.Cards(
					ct.Bind(&inPlay, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits("Niffle"))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Traits("Beast"))),
				),
			},
		})

		h.P1.Play(TroopCall)

		h.Expect(buried).At(ct.Hand)
		h.Expect(inPlay).At(ct.Hand)
		h.Expect(other).At(ct.PlayArea)
	})
}
