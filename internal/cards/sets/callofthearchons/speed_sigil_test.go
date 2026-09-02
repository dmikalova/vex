package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Speed Sigil
func TestSpeedSigil(t *testing.T) {
	shadows := ct.OfHouse(card.House.Shadows)

	t.Run("the first creature played enters play ready", func(t *testing.T) {
		var thief ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(SpeedSigil),
				Hand:   ct.Cards(ct.Bind(&thief, ct.Creature(shadows))),
			},
		})

		h.P1.Play(thief)
		h.Expect(thief).Ready()
	})

	t.Run("later creatures enter play exhausted", func(t *testing.T) {
		var first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(SpeedSigil),
				Hand: ct.Cards(
					ct.Bind(&first, ct.Creature(shadows)),
					ct.Bind(&second, ct.Creature(shadows)),
				),
			},
		})

		h.P1.Play(first)
		h.P1.Play(second)
		h.Expect(first).Ready()
		h.Expect(second).Exhausted()
	})

	t.Run("the charge resets each turn", func(t *testing.T) {
		var first, third ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(SpeedSigil),
				Hand: ct.Cards(
					ct.Bind(&first, ct.Creature(shadows)),
					ct.Bind(&third, ct.Creature(shadows)),
				),
			},
			P2: ct.Side{House: card.House.Shadows},
		})

		h.P1.Play(first)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Shadows)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Shadows)
		h.P1.Play(third)
		h.Expect(third).Ready()
	})
}
