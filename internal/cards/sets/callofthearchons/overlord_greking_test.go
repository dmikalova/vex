package callofthearchons

import (
	"slices"
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Overlord Greking
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	After a creature is destroyed in a fight with Overlord Greking, put it into play under your control.
func TestOverlordGreking(t *testing.T) {
	t.Run(
		"reanimates an enemy creature destroyed in a fight with it, under your control",
		func(t *testing.T) {
			var foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Dis,
					InPlay: ct.Cards(OverlordGreking),
				},
				P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
			})

			h.P1.Fight(OverlordGreking, foe)

			if !slices.Contains(h.Game().Battleline(0), foe.ID()) {
				t.Error("the destroyed enemy should be reanimated under P1's control")
			}
			h.Expect(foe).At(ct.PlayArea)
		},
	)
}
