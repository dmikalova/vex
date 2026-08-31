package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Skeleton Key
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: A friendly creature captures 1 Æmber from your opponent.
func TestSkeletonKey(t *testing.T) {
	t.Run("a friendly creature captures 1 Æmber from the opponent", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					SkeletonKey,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Shadows))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(SkeletonKey)

		h.Expect(ally).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
