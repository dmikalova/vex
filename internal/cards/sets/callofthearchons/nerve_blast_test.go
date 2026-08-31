package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nerve Blast
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber -> deal 2 damage to a creature.
func TestNerveBlast(t *testing.T) {
	t.Run("steals 1 Æmber and, if it does, deals 2 damage", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(NerveBlast)},
			P2: ct.Side{
				Amber: 3,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Play(NerveBlast)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
		h.Expect(foe).Damage(2)
	})

	t.Run("does nothing when there is no Æmber to steal", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(NerveBlast)},
			P2: ct.Side{
				Amber: 0,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Play(NerveBlast)

		h.Expect(foe).Damage(0) // the gate did not fire
	})
}
