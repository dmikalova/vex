package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tolas
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Imp
//
//	Elusive.
//	Each creature gains, "Destroyed: Your opponent gains 1 Æmber."
func TestTolas(t *testing.T) {
	t.Run("a destroyed creature's opponent gains 1 Æmber", func(t *testing.T) {
		var attacker, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Tolas,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Fight(attacker, foe)

		h.Expect(foe).At(ct.Discard)
		h.P1.ExpectAmber(1) // the destroyed creature's opponent
	})
}
