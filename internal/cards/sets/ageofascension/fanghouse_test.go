package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fanghouse
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Assault 3, Hazardous 3.
func TestFanghouse(t *testing.T) {
	t.Run("deals assault damage before combat", func(t *testing.T) {
		var fanghouse, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&fanghouse, Fanghouse)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Fight(fanghouse, foe)

		h.Expect(foe).Damage(6)
	})

	t.Run("deals hazardous damage to an attacker", func(t *testing.T) {
		var attacker, fanghouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&fanghouse, Fanghouse))},
		})

		h.P1.Fight(attacker, fanghouse)

		h.Expect(attacker).At(ct.Discard)
	})
}
