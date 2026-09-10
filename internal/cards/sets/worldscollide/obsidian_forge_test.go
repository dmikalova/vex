package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Obsidian Forge
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Destroy any number of friendly creatures, then forge a key at +6 Æmber current cost, reduced by 1 Æmber for each creature destroyed this way -> purge Obsidian Forge.
func TestObsidianForge(t *testing.T) {
	t.Run(
		"sacrifices creatures, forges the reduced-cost key, and purges itself",
		func(t *testing.T) {
			var a, b ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					Amber: 10,
					InPlay: ct.Cards(
						ObsidianForge,
						ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Dis))),
						ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Dis))),
					),
				},
			})

			h.P1.UseAction(ObsidianForge)
			h.P1.ClickCard(a)
			h.P1.ClickCard(b)

			h.Expect(a).At(ct.Discard)
			h.Expect(b).At(ct.Discard)
			h.Expect(ObsidianForge).At(ct.Purge)
			h.P1.ExpectKeys(1)
			// Cost was 6 + 6 - 2 sacrificed = 10, spending the whole pool.
			h.P1.ExpectAmber(0)
		},
	)

	t.Run("an unaffordable forge leaves Obsidian Forge in play", func(t *testing.T) {
		var a ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Amber: 10,
				InPlay: ct.Cards(
					ObsidianForge,
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.UseAction(ObsidianForge)
		h.P1.ClickCard(a)

		h.Expect(a).At(ct.Discard)
		h.Expect(ObsidianForge).At(ct.PlayArea)
		h.P1.ExpectKeys(0)
		h.P1.ExpectAmber(10)
	})
}
