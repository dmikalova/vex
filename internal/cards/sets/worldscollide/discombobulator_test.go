package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Discombobulator
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Your Æmber cannot be stolen."
func TestDiscombobulator(t *testing.T) {
	t.Run("protects its controller's Æmber while attached", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature()), Discombobulator),
				),
				Amber: 3,
			},
		})

		if !h.Game().AemberProtected(0) {
			t.Error("Discombobulator should make its host's controller's Æmber unstealable")
		}
	})
}
