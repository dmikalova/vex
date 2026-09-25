package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Munchling
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant
//
//	Skirmish.
//	Fight: You may discard a Logos card from your hand or archives -> gain 1 Æmber.
func TestMunchling(t *testing.T) {
	t.Run("discards a Logos card from hand and gains 1 Æmber", func(t *testing.T) {
		var munchling, logos, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&munchling, Munchling)),
				Hand:   ct.Cards(ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Fight(munchling, foe)
		h.P1.ClickOption("Yes") // the sole Logos card is discarded automatically

		h.Expect(logos).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})

	t.Run("may discard a Logos card drawn from the archives", func(t *testing.T) {
		var munchling, inHand, inArchives, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:    card.House.Logos,
				InPlay:   ct.Cards(ct.Bind(&munchling, Munchling)),
				Hand:     ct.Cards(ct.Bind(&inHand, ct.Creature(ct.OfHouse(card.House.Logos)))),
				Archives: ct.Cards(ct.Bind(&inArchives, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Fight(munchling, foe)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(inArchives) // pick the archived card over the one in hand

		h.Expect(inArchives).At(ct.Discard)
		h.Expect(inHand).At(ct.Hand)
		h.P1.ExpectAmber(1)
	})

	t.Run("declining discards nothing and gains nothing", func(t *testing.T) {
		var munchling, logos, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&munchling, Munchling)),
				Hand:   ct.Cards(ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Fight(munchling, foe)
		h.P1.ClickOption("No")

		h.Expect(logos).At(ct.Hand)
		h.P1.ExpectAmber(0)
	})
}
