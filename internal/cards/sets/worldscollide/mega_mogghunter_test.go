package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mega Mogghunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  8
//	Traits: Giant
//
//	Fight: Deal 2 damage to a flank creature.
func TestMegaMogghunter(t *testing.T) {
	t.Run("deals 2 damage to a flank creature when it fights", func(t *testing.T) {
		var bruiser, squishy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(MegaMogghunter),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&bruiser, ct.Creature(ct.Power(3), ct.Armor(8))),
				ct.Bind(&squishy, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Fight(MegaMogghunter, bruiser)
		h.P1.ClickCard(squishy)

		h.Expect(squishy).Damage(2)
	})
}
