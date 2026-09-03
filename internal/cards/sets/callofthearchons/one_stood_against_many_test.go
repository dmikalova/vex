package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestOneStoodAgainstMany(t *testing.T) {
	t.Run("fights three times, never the same enemy twice", func(t *testing.T) {
		var hero, foe1, foe2, foe3 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(OneStoodAgainstMany),
				InPlay: ct.Cards(
					ct.Bind(&hero, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(10))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe1, ct.Creature(ct.Power(1))),
					ct.Bind(&foe2, ct.Creature(ct.Power(1))),
					ct.Bind(&foe3, ct.Creature(ct.Power(1))),
				),
			},
		})

		h.P1.Play(OneStoodAgainstMany)
		// The lone friendly creature is automatic; each fight narrows the enemies
		// left, so only the first two need a click.
		h.P1.ClickCard(foe1)
		h.P1.ClickCard(foe2)

		h.Expect(foe1).At(ct.Discard)
		h.Expect(foe2).At(ct.Discard)
		h.Expect(foe3).At(ct.Discard)
		// One point of damage back from each of the three fights, and the last
		// fight leaves the hero exhausted.
		h.Expect(hero).Damage(3).Exhausted()
	})

	t.Run("stops when there is no enemy left to fight", func(t *testing.T) {
		var hero, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(OneStoodAgainstMany),
				InPlay: ct.Cards(
					ct.Bind(&hero, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(10))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1))))},
		})

		h.P1.Play(OneStoodAgainstMany)

		h.Expect(foe).At(ct.Discard)
		h.Expect(hero).Damage(1)
	})
}
