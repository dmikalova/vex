package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Floomf
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast • Cat
//
//	Skirmish.
//	Fight: Give a Beast creature two +1 power counters.
func TestFloomf(t *testing.T) {
	t.Run("gives a chosen Beast creature two +1 power counters when it fights", func(t *testing.T) {
		var beast, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					Floomf,
					ct.Bind(
						&beast,
						ct.Creature(
							ct.OfHouse(card.House.Untamed),
							ct.Power(3),
							ct.Traits(card.Traits.Beast),
						),
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
			)},
		})

		h.P1.Fight(Floomf, enemy)
		// Floomf is also a Beast, so the "a Beast creature" target prompts a choice.
		h.P1.ClickCard(beast)

		h.Expect(beast).Power(5)
	})
}
