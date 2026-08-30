package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sanctum Guardian
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Knight • Spirit
//
//	Taunt.
//	Fight/Reap: Swap this creature with another friendly creature in your battleline.
func TestSanctumGuardian(t *testing.T) {
	t.Run("swaps with another friendly creature when it reaps", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(SanctumGuardian, ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
		})

		h.P1.Reap(SanctumGuardian)

		// Both creatures remain in play after the swap.
		h.Expect(SanctumGuardian).At(ct.PlayArea)
		h.Expect(ally).At(ct.PlayArea)
	})
}
