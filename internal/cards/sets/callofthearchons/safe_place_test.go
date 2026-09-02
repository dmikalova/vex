package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Safe Place
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	You may spend Æmber on Safe Place when forging keys.
//	Action: Move 1 Æmber from your pool to Safe Place.
func TestSafePlace(t *testing.T) {
	t.Run("banks 1 Æmber from the pool", func(t *testing.T) {
		var place ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&place, SafePlace)),
				Amber:  3,
			},
		})

		h.P1.UseAction(place)

		h.P1.ExpectAmber(2)
		h.Expect(place).AmberOn(1)
	})

	t.Run("banks nothing from an empty pool", func(t *testing.T) {
		var place ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&place, SafePlace)),
			},
		})

		h.P1.UseAction(place)

		h.Expect(place).AmberOn(0)
	})

	t.Run("its banked Æmber pays for a key", func(t *testing.T) {
		var place ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&place, SafePlace)),
				Amber:  6,
			},
		})

		h.P1.UseAction(place)
		h.P1.ExpectAmber(5)
		h.Expect(place).AmberOn(1)

		h.P1.EndTurn()
		h.P2.EndTurn()

		// 5 in the pool plus the 1 banked covers the 6-Æmber key.
		h.P1.ExpectKeys(1)
		h.P1.ExpectAmber(0)
		h.Expect(place).AmberOn(0)
	})
}
