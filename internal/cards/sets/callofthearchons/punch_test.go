package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Punch
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 3 damage to a creature.
func TestPunch(t *testing.T) {
	t.Run("deals 3 damage to a chosen creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Punch)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)))),
			},
		})

		h.P1.Play(Punch)

		h.Expect(foe).At(ct.PlayArea).Damage(3)
	})
}
