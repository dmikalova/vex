package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Vaultkeeper
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Knight • Spirit
//
//	Your Æmber cannot be stolen.
func TestTheVaultkeeper(t *testing.T) {
	t.Run("protects its controller's Æmber from being stolen", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(TheVaultkeeper), Amber: 3},
		})

		if !h.Game().AemberProtected(0) {
			t.Error("The Vaultkeeper should make its controller's Æmber unstealable")
		}
	})
}
