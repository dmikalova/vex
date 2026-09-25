package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mogghunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Fight: Deal 2 damage to a flank creature.
func TestMogghunter(t *testing.T) {
	t.Run("deals 2 damage to a flank creature when it fights", func(t *testing.T) {
		var bruiser, squishy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Mogghunter),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&bruiser, ct.Creature(ct.Power(3), ct.Armor(8))),
				ct.Bind(&squishy, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Fight(Mogghunter, bruiser)
		h.P1.ClickCard(squishy)

		h.Expect(squishy).Damage(2)
		h.Expect(Mogghunter).At(ct.PlayArea)
	})
}
