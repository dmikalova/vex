package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dance of Doom
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature - destroy each creature with the same power as the chosen creature.
func TestDanceOfDoom(t *testing.T) {
	t.Run("destroys every creature sharing the chosen creature's power", func(t *testing.T) {
		var weak, strong, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
					ct.Bind(&strong, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(6))),
				),
				Hand: ct.Cards(DanceOfDoom),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})

		h.P1.Play(DanceOfDoom)
		h.P1.ExpectPrompt("Choose a creature").Source("Dance of Doom")
		h.P1.ClickCard(weak)

		h.Expect(strong).At(ct.PlayArea) // power 6, unmatched
		h.Expect(weak).At(ct.Discard)    // the chosen power-3 creature
		h.Expect(enemy).At(ct.Discard)   // power 3, matched across both battlelines
	})
}
