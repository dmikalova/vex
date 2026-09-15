package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Floor is Lava
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	At the start of your turn, deal 1 damage to a friendly creature, and deal 1 damage to an enemy creature.
func TestTheFloorIsLava(t *testing.T) {
	t.Run(
		"deals 1 damage to a friendly and an enemy creature at the start of your turn",
		func(t *testing.T) {
			var friend, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Brobnar,
					InPlay: ct.Cards(TheFloorIsLava, ct.Bind(&friend, ct.Creature(ct.Power(5)))),
				},
				P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(5))))},
			})

			h.P1.EndTurn() // to P2
			h.P2.EndTurn() // back to P1: start-of-turn trigger fires

			h.Expect(friend).Damage(1)
			h.Expect(enemy).Damage(1)
		},
	)
}
