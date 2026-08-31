package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Strange Gizmo
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After you forge a key, destroy each creature and each artifact.
func TestStrangeGizmo(t *testing.T) {
	t.Run(
		"destroys each creature and artifact after its controller forges a key",
		func(t *testing.T) {
			var ally, foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Logos,
					Amber: 6,
					InPlay: ct.Cards(
						StrangeGizmo,
						ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					),
				},
			})

			h.P1.EndTurn() // P1's turn ends, P2 begins
			h.P2.EndTurn() // P1 begins and forges a key with 6 Æmber, firing Strange Gizmo

			h.Expect(ally).At(ct.Discard)
			h.Expect(foe).At(ct.Discard)
			h.Expect(StrangeGizmo).At(ct.Discard)
		},
	)
}
