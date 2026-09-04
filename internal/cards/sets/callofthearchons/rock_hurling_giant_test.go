package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rock-Hurling Giant
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant
//
//	After you discard a card from your hand, if it is a Brobnar card, you may deal 4 damage to a creature.
func TestRockHurlingGiant(t *testing.T) {
	t.Run("may deal 4 damage when a Brobnar card is discarded", func(t *testing.T) {
		var giant, fodder, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&giant, RockHurlingGiant)),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Discard(fodder)
		h.P1.ClickCard(enemy)

		h.Expect(enemy).Damage(4)
	})

	t.Run("stays quiet for an off-house discard", func(t *testing.T) {
		var giant, fodder, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&giant, RockHurlingGiant)),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Discard(fodder)

		h.Expect(enemy).Damage(0)
	})
}
