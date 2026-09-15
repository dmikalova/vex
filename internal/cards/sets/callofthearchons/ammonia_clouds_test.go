package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ammonia Clouds
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 3 damage to each creature.
func TestAmmoniaClouds(t *testing.T) {
	t.Run("deals 3 damage to each creature", func(t *testing.T) {
		var toughAlly, weakAlly, toughFoe, weakFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(AmmoniaClouds),
				InPlay: ct.Cards(
					ct.Bind(&toughAlly, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
					ct.Bind(&weakAlly, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&toughFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
					ct.Bind(&weakFoe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
				),
			},
		})

		h.P1.Play(AmmoniaClouds)

		// Weak creatures (power <= 3) die; the tough bodies survive marked with 3.
		h.Expect(weakAlly).At(ct.Discard)
		h.Expect(weakFoe).At(ct.Discard)
		h.Expect(toughAlly).At(ct.PlayArea).Damage(3)
		h.Expect(toughFoe).At(ct.PlayArea).Damage(3)
	})
}
