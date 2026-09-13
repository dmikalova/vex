package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hit and Run
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2 damage to a Creature. Put a friendly Creature into its owner's hand.
func TestHitAndRun(t *testing.T) {
	t.Run(
		"deals 2 damage to a creature and returns a friendly creature to hand",
		func(t *testing.T) {
			var foe, ally ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					Hand:   ct.Cards(HitAndRun),
					InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Shadows)))),
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4)))),
				},
			})

			h.P1.Play(HitAndRun)
			h.P1.ClickCard(foe)

			h.Expect(foe).At(ct.PlayArea).Damage(2)
			h.Expect(ally).At(ct.Hand)
		},
	)
}
