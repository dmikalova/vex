package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Honorable Claim
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each friendly Knight trait creature captures 1 Æmber from your opponent.
func TestHonorableClaim(t *testing.T) {
	t.Run("each friendly Knight captures 1 Æmber", func(t *testing.T) {
		var knight, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(HonorableClaim),
				InPlay: ct.Cards(
					ct.Bind(
						&knight,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Knight)),
					),
					ct.Bind(
						&other,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Traits(card.Traits.Human)),
					),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(HonorableClaim)

		h.Expect(knight).AmberOn(1)
		h.Expect(other).AmberOn(0)
		h.P2.ExpectAmber(2)
	})
}
