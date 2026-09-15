package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Forge Compiler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	After your opponent forges a key, destroy Forge Compiler, and ward each friendly creature.
func TestForgeCompiler(t *testing.T) {
	t.Run(
		"destroys itself and wards friendly creatures after the opponent forges",
		func(t *testing.T) {
			var ally ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Logos,
					InPlay: ct.Cards(
						ForgeCompiler,
						ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Logos))),
					),
				},
			})

			g := h.Game()
			g.State.Aember[1] = engine.KeyCost
			g.StartTurn(1)

			h.Expect(ForgeCompiler).At(ct.Discard)
			if !g.Warded(ally.ID()) {
				t.Error("friendly creature should be warded after the opponent forges")
			}
		},
	)
}
