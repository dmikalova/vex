package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Savage Clash
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each enemy creature except the most powerful enemy creature and each friendly creature except the least powerful friendly creature.
func TestSavageClash(t *testing.T) {
	t.Run("spares the strongest enemy and the weakest ally", func(t *testing.T) {
		var strongAlly, weakAlly, strongFoe, weakFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(SavageClash),
				InPlay: ct.Cards(
					ct.Bind(&strongAlly, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(5))),
					ct.Bind(&weakAlly, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(2))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&strongFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
					ct.Bind(&weakFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(SavageClash)

		h.Expect(strongFoe).At(ct.PlayArea) // most powerful enemy survives
		h.Expect(weakAlly).At(ct.PlayArea)  // least powerful friendly survives
		h.Expect(weakFoe).At(ct.Discard)
		h.Expect(strongAlly).At(ct.Discard)
	})
}
