package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Opal Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Spirit • Knight
//
//	Play: Destroy each creature with even power.
func TestOpalKnight(t *testing.T) {
	t.Run("destroys each creature with even power", func(t *testing.T) {
		var even, odd ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(OpalKnight)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&even, ct.Creature(ct.Power(4))),
				ct.Bind(&odd, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Play(OpalKnight)

		h.Expect(even).At(ct.Discard)
		h.Expect(odd).At(ct.PlayArea)
	})
}
