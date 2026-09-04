package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestShadowSelf(t *testing.T) {
	t.Run("shields a non-Specter neighbor and deals no fight damage", func(t *testing.T) {
		var shadow, ward, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(
				ct.Bind(&ward, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(2))),
				ct.Bind(&shadow, ShadowSelf),
			)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		// Shadow Self's 9 power would kill the 3-power enemy, but it deals none.
		h.P1.Fight(shadow, foe)
		h.Expect(foe).At(ct.PlayArea).Damage(0)
		// The enemy still fights back, and that damage lands on Shadow Self anyway.
		h.Expect(shadow).At(ct.PlayArea).Damage(3)

		// Damage aimed at the weak neighbor lands on Shadow Self instead, so the
		// neighbor that a 3-power fight would have killed is untouched.
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars)
		h.P2.Fight(foe, ward)
		h.Expect(ward).At(ct.PlayArea).Damage(0)
		h.Expect(shadow).At(ct.PlayArea).Damage(6)
	})

	t.Run("does not shield a Specter neighbor", func(t *testing.T) {
		var shadow, specter, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(
				ct.Bind(&specter, ct.Creature(
					ct.OfHouse(card.House.Shadows), ct.Power(5), ct.Traits(card.Traits.Specter))),
				ct.Bind(&shadow, ShadowSelf),
			)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars)
		h.P2.Fight(foe, specter)
		h.Expect(specter).At(ct.PlayArea).Damage(3)
		h.Expect(shadow).At(ct.PlayArea).Damage(0)
	})
}
