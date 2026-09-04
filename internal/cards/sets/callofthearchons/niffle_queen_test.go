package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Niffle Queen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Beast • Niffle
//
//	Each other friendly Beast trait creature gains +1 power.
//	Each other friendly Niffle trait creature gains +1 power.
func TestNiffleQueen(t *testing.T) {
	t.Run("boosts other friendly Beasts and Niffles but not itself or enemies", func(t *testing.T) {
		var queen, niffleBeast, beast, plain, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&queen, NiffleQueen),
					ct.Bind(&niffleBeast, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Power(3),
						ct.Traits(card.Traits.Beast, card.Traits.Niffle),
					)),
					ct.Bind(&beast, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Power(3),
						ct.Traits(card.Traits.Beast),
					)),
					ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(
						ct.OfHouse(card.House.Untamed),
						ct.Power(3),
						ct.Traits(card.Traits.Beast),
					)),
				),
			},
		})

		h.Expect(queen).Power(6)
		h.Expect(niffleBeast).Power(5) // both constants apply
		h.Expect(beast).Power(4)
		h.Expect(plain).Power(3)
		h.Expect(foe).Power(3)
	})
}
