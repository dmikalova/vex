package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ogopogo
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed fighting Ogopogo, you may deal 2 damage to a creature.
func TestOgopogo(t *testing.T) {
	t.Run("may deal 2 damage after destroying a creature in a fight", func(t *testing.T) {
		var ogopogo, foe, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&ogopogo, Ogopogo)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(1))),
					ct.Bind(&other, ct.Creature(ct.Power(9))),
				),
			},
		})

		h.P1.Fight(ogopogo, foe)
		h.P1.ClickCard(other)

		h.Expect(foe).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea).Damage(2)
	})

	t.Run("may decline the damage", func(t *testing.T) {
		var ogopogo, foe, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&ogopogo, Ogopogo)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(1))),
					ct.Bind(&other, ct.Creature(ct.Power(9))),
				),
			},
		})

		h.P1.Fight(ogopogo, foe)
		h.P1.ClickDone()

		h.Expect(foe).At(ct.Discard)
		h.Expect(other).At(ct.PlayArea).Damage(0)
	})
}
