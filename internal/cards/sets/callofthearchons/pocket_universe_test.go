package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pocket Universe
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	You may spend Æmber on Pocket Universe when forging keys.
//	Action: Move 1 Æmber from your pool to Pocket Universe.
func TestPocketUniverse(t *testing.T) {
	t.Run("banks 1 Æmber from the pool", func(t *testing.T) {
		var universe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&universe, PocketUniverse)),
				Amber:  2,
			},
		})

		h.P1.UseAction(universe)

		h.P1.ExpectAmber(1)
		h.Expect(universe).AmberOn(1)
	})
}
