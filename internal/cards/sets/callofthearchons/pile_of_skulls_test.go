package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pile of Skulls
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	After an enemy creature is destroyed during your turn, a friendly creature captures 1 Æmber from your opponent.
func TestPileOfSkulls(t *testing.T) {
	t.Run(
		"a friendly creature captures 1 when an enemy creature is destroyed on your turn",
		func(t *testing.T) {
			var fighter, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(
					PileOfSkulls,
					ct.Bind(&fighter, ct.Creature(ct.Power(5))),
				)},
				P2: ct.Side{Amber: 3, InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(2))))},
			})

			h.P1.Fight(fighter, enemy)

			h.Expect(fighter).AmberOn(1) // the sole friendly creature captures
			h.P2.ExpectAmber(2)          // 1 taken from the opponent's pool
		},
	)
}
