package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Blast from the Past
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Exalt a friendly creature. Archive a Saurian creature from your discard pile. Deal damage equal to its power to an enemy creature.
func TestBlastFromThePast(t *testing.T) {
	t.Run("exalts, archives a Saurian, and deals its power to an enemy", func(t *testing.T) {
		var friendly, saurian, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(BlastFromThePast),
				InPlay: ct.Cards(ct.Bind(&friendly, ct.Creature(ct.Power(4)))),
				Discard: ct.Cards(
					ct.Bind(&saurian, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Play(BlastFromThePast)

		h.Expect(friendly).AmberOn(1)
		h.Expect(saurian).At(ct.Archives)
		h.Expect(enemy).Damage(3)
	})

	t.Run("no Saurian in the discard pile deals no damage", func(t *testing.T) {
		var friendly, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(BlastFromThePast),
				InPlay: ct.Cards(ct.Bind(&friendly, ct.Creature(ct.Power(4)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Play(BlastFromThePast)

		h.Expect(friendly).AmberOn(1)
		h.Expect(enemy).Damage(0)
	})
}
