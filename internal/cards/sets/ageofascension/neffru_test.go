package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Neffru
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	After a creature is destroyed, its owner gains 1 Æmber.
func TestNeffru(t *testing.T) {
	t.Run("a destroyed creature's owner gains 1 aember", func(t *testing.T) {
		var attacker, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Neffru,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(2)))),
			},
		})

		h.P1.Fight(attacker, enemy)

		h.Expect(enemy).At(ct.Discard)
		h.P2.ExpectAmber(1)
	})

	t.Run("each creature destroyed together pays its own owner", func(t *testing.T) {
		var attacker, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					Neffru,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Fight(attacker, enemy)

		h.Expect(attacker).At(ct.Discard)
		h.Expect(enemy).At(ct.Discard)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("neffru pays nothing when it is destroyed in the batch", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(Neffru),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Fight(Neffru, enemy)

		h.Expect(Neffru).At(ct.Discard)
		h.Expect(enemy).At(ct.Discard)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(0)
	})
}
