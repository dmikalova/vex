package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// A. Vinda
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Elf • Thief
//
//	Reap: Deal 1 damage to a creature. If this damage destroys that creature, your opponent discards a random card from their hand.
func TestAVinda(t *testing.T) {
	t.Run("opponent discards when the damage destroys the target", func(t *testing.T) {
		var vinda, foe, held ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&vinda, AVinda)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(1)))),
				Hand:   ct.Cards(ct.Bind(&held, ct.Creature())),
			},
		})

		h.P1.Reap(vinda)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(held).At(ct.Discard)
	})

	t.Run("no discard when the target survives", func(t *testing.T) {
		var vinda, foe, held ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&vinda, AVinda)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4)))),
				Hand:   ct.Cards(ct.Bind(&held, ct.Creature())),
			},
		})

		h.P1.Reap(vinda)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.PlayArea)
		h.Expect(held).At(ct.Hand)
	})
}
