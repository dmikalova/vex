package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Look Over There!
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2 damage to a creature. If it is not destroyed, steal 1 Æmber.
func TestLookOverThere(t *testing.T) {
	t.Run("steals 1 Æmber when the creature survives the damage", func(t *testing.T) {
		var tough ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(LookOverThere),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&tough, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
				Amber: 2,
			},
		})

		h.P1.Play(LookOverThere)

		h.Expect(tough).At(ct.PlayArea).Damage(2)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})

	t.Run("does not steal when the damage destroys the creature", func(t *testing.T) {
		var frail ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(LookOverThere),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&frail, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(1))),
				),
				Amber: 2,
			},
		})

		h.P1.Play(LookOverThere)

		h.Expect(frail).At(ct.Discard)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(2)
	})
}
