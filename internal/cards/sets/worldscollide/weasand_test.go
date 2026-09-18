package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Weasand
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Thief
//
//	Deploy, Elusive.
//	If Weasand is on a flank, destroy Weasand.
//	After a player forges a key, gain 2 Æmber.
func TestWeasand(t *testing.T) {
	t.Run("is destroyed while it holds a flank", func(t *testing.T) {
		var weasand, filler ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&weasand, Weasand)),
				Hand: ct.Cards(ct.Bind(&filler,
					ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3)))),
			},
			P2: ct.Side{},
		})

		// Playing another creature settles the board: Weasand still holds a
		// flank, so its condition destroys it.
		h.P1.Play(filler)

		h.Expect(weasand).At(ct.Discard)
	})

	t.Run(
		"survives shielded by neighbors, then gains 2 Æmber when a key is forged",
		func(t *testing.T) {
			var weasand ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Shadows,
					InPlay: ct.Cards(
						ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3)),
						ct.Bind(&weasand, Weasand),
						ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3)),
					),
				},
				P2: ct.Side{},
			})
			g := h.Game()
			g.State.ForgeCanonicalKeys(0, 0)
			g.State.Aember[0] = engine.KeyCost

			h.Expect(weasand).At(ct.PlayArea)

			h.P1.EndTurn() // to P2
			h.P2.EndTurn() // back to P1: forge phase forges a key, Weasand gains 2

			h.P1.ExpectKeys(1)
			h.P1.ExpectAmber(2)
			h.Expect(weasand).At(ct.PlayArea)
		},
	)
}
