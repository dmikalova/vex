package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Orbital Bombardment
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Reveal any number of Mars cards from your hand. For each card revealed this way, deal 2 damage to a creature.
func TestOrbitalBombardment(t *testing.T) {
	t.Run("deals 2 damage per Mars card revealed", func(t *testing.T) {
		var bombardment, mars, other, first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					ct.Bind(&bombardment, OrbitalBombardment),
					ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&first, ct.Creature(ct.Power(9))),
					ct.Bind(&second, ct.Creature(ct.Power(9))),
				),
			},
		})

		h.P1.Play(bombardment)
		// Only Mars cards may be revealed, and the reveal is "any number".
		h.P1.ClickCard(mars)
		h.P1.ClickCard(first)

		h.Expect(first).Damage(2)
		h.Expect(second).Damage(0)
		h.Expect(other).At(ct.Hand)
	})

	t.Run("damages a different creature for each reveal", func(t *testing.T) {
		var bombardment, marsOne, marsTwo, first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					ct.Bind(&bombardment, OrbitalBombardment),
					ct.Bind(&marsOne, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Bind(&marsTwo, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&first, ct.Creature(ct.Power(9))),
					ct.Bind(&second, ct.Creature(ct.Power(9))),
				),
			},
		})

		h.P1.Play(bombardment)
		h.P1.ClickCard(marsOne)
		h.P1.ClickCard(marsTwo)
		h.P1.ClickCard(first)
		h.P1.ClickCard(second)

		h.Expect(first).Damage(2)
		h.Expect(second).Damage(2)
	})

	t.Run("deals nothing when nothing is revealed", func(t *testing.T) {
		var bombardment, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&bombardment, OrbitalBombardment)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(bombardment)

		h.Expect(enemy).Damage(0)
	})
}
