package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// J. Vinda
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: Deal 1 damage to a creature. If this damage destroys that creature, steal 1 Æmber.
func TestJVinda(t *testing.T) {
	t.Run("reaps to destroy a creature and steal 1 Æmber", func(t *testing.T) {
		var vinda, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&vinda, JVinda)),
			},
			P2: ct.Side{
				Amber:  1,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Reap(vinda)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.Discard)
		h.P2.ExpectAmber(0)
	})
}
