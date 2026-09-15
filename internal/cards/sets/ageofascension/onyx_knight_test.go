package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Onyx Knight
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon • Knight
//
//	Play: Destroy each creature with odd power.
func TestOnyxKnight(t *testing.T) {
	t.Run("destroys each creature with odd power", func(t *testing.T) {
		var odd, even ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(OnyxKnight)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&odd, ct.Creature(ct.Power(3))),
				ct.Bind(&even, ct.Creature(ct.Power(4))),
			)},
		})

		h.P1.Play(OnyxKnight)

		h.Expect(odd).At(ct.Discard)
		h.Expect(even).At(ct.PlayArea)
	})
}
