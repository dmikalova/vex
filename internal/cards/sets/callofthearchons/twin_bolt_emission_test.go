package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Twin Bolt Emission
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to a creature and deal 2 damage to a different creature.
func TestTwinBoltEmission(t *testing.T) {
	t.Run("deals 2 to a creature and 2 to a different creature", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(TwinBoltEmission)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.Power(5))),
				ct.Bind(&b, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Play(TwinBoltEmission)
		h.P1.ClickCard(a) // first creature; the different creature is the sole remaining one

		h.Expect(a).Damage(2)
		h.Expect(b).Damage(2)
	})
}
