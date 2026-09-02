package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mothership Support
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each friendly ready Mars creature, deal 2 damage to a creature.
func TestMothershipSupport(t *testing.T) {
	t.Run("deals 2 damage per friendly ready Mars creature", func(t *testing.T) {
		var support, exhausted, first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&support, MothershipSupport)),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
					ct.Bind(&exhausted, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&first, ct.Creature(ct.Power(9))),
					ct.Bind(&second, ct.Creature(ct.Power(9))),
				),
			},
		})
		exhausted.Exhaust()

		h.P1.Play(support)
		// Two ready Mars creatures, so two separate 2-damage choices.
		h.P1.ClickCard(first)
		h.P1.ClickCard(second)

		h.Expect(first).Damage(2)
		h.Expect(second).Damage(2)
	})

	t.Run("deals nothing with no ready Mars creature", func(t *testing.T) {
		var support, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&support, MothershipSupport)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(support)

		h.Expect(enemy).Damage(0)
	})
}
