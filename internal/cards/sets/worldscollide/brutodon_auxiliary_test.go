package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Brutodon Auxiliary
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Beast
//
//	Taunt, Hazardous 2.
func TestBrutodonAuxiliary(t *testing.T) {
	t.Run("deals 2 hazardous damage to an attacker before combat", func(t *testing.T) {
		var brutodon, attacker ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(ct.Power(20)))),
			},
			P2: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&brutodon, BrutodonAuxiliary)),
			},
		})
		attacker.Ready()

		h.P1.Fight(attacker, brutodon)

		h.Expect(brutodon).At(ct.Discard)
		// 6 fight damage from Brutodon's power plus 2 from Hazardous.
		h.Expect(attacker).At(ct.PlayArea).Damage(8)
	})
}
