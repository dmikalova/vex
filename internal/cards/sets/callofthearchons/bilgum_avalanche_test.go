package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Bilgum Avalanche
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Giant
//
//	After you forge a key, deal 2 damage to each enemy creature.
func TestBilgumAvalanche(t *testing.T) {
	t.Run(
		"deals 2 damage to each enemy creature when its controller forges a key",
		func(t *testing.T) {
			var foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Brobnar,
					InPlay: ct.Cards(BilgumAvalanche),
					Amber:  engine.KeyCost,
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
					),
				},
			})

			h.Game().StartTurn(0) // forge a key with the seeded Æmber

			h.Expect(foe).Damage(2)
		},
	)
}
