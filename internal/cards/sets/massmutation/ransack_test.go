package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ransack
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Steal 1 Æmber. Discard the top card of your deck. -> if the discarded card is a Shadows card, repeat this effect.
func TestRansack(t *testing.T) {
	t.Run("steals and repeats while the discarded card is a Shadows card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Ransack),
				Deck: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Shadows)),
					ct.Creature(ct.OfHouse(card.House.Untamed)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(Ransack)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(3)
	})

	t.Run("stops immediately when the discarded card is not a Shadows card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Ransack),
				Deck: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Untamed)),
					ct.Creature(ct.OfHouse(card.House.Shadows)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(Ransack)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(4)
	})
}
