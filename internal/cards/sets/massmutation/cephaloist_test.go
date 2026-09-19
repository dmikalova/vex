package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cephaloist
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	While you have 4 Æmber or more, your Æmber cannot be stolen.
func TestCephaloist(t *testing.T) {
	t.Run("protects its controller's Æmber while their pool is at least 4", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Cephaloist),
				Amber:  4,
			},
		})

		if !h.Game().AemberProtected(0) {
			t.Error("Cephaloist should make its controller's Æmber unstealable at 4 Æmber")
		}
	})

	t.Run("does not protect while the pool is below 4", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(Cephaloist),
				Amber:  3,
			},
		})

		if h.Game().AemberProtected(0) {
			t.Error("Cephaloist should not protect Æmber below 4 Æmber")
		}
	})
}
