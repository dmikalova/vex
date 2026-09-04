package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// "John Smyth"
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Agent • Martian
//
//	Elusive.
//	Fight/Reap: Ready a non-Agent trait Mars creature.
func TestJohnSmyth(t *testing.T) {
	t.Run("readies a non-Agent Mars creature on reap", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					JohnSmyth,
					ct.Bind(
						&ally,
						ct.Creature(
							ct.OfHouse(card.House.Mars),
							ct.Traits(card.Traits.Martian),
							ct.Power(3),
						),
					),
				),
			},
		})
		ally.Exhaust()

		h.P1.Reap(JohnSmyth)

		h.Expect(ally).Ready()
	})

	t.Run("cannot ready an Agent or a non-Mars creature", func(t *testing.T) {
		var agent, offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					JohnSmyth,
					ct.Bind(
						&agent,
						ct.Creature(
							ct.OfHouse(card.House.Mars),
							ct.Traits(card.Traits.Agent),
							ct.Power(3),
						),
					),
					ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})
		agent.Exhaust()
		offHouse.Exhaust()

		h.P1.Reap(JohnSmyth)

		h.Expect(agent).Exhausted()
		h.Expect(offHouse).Exhausted()
	})
}
