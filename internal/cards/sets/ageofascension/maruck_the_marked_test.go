package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Maruck the Marked
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Spirit • Knight
//
//	After Maruck the Marked prevents damage with its armor, for each damage just prevented, Maruck the Marked captures 1 Æmber from your opponent.
func TestMaruckTheMarked(t *testing.T) {
	t.Run("captures 1 Æmber for each damage its armor prevents", func(t *testing.T) {
		var maruck, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&maruck, MaruckTheMarked)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
				Amber:  4,
			},
		})

		h.P1.Fight(maruck, foe)

		// Maruck's 1 armor absorbs 1 of the 3 damage dealt back, so it captures 1.
		h.Expect(maruck).AmberOn(1)
		h.Expect(maruck).Damage(2)
	})

	t.Run("does not fire when no armor is spent", func(t *testing.T) {
		var maruck, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ct.Bind(&maruck, MaruckTheMarked)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(0)))),
				Amber:  4,
			},
		})

		h.P1.Fight(maruck, foe)

		h.Expect(maruck).AmberOn(0)
	})
}
