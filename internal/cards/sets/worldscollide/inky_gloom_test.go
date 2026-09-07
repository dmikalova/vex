package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Inky Gloom
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent cannot use creatures to reap during their next turn.
func TestInkyGloom(t *testing.T) {
	t.Run(
		"bars the opponent from reaping on their next turn, fighting stays open",
		func(t *testing.T) {
			var foe, ally ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					Hand:   ct.Cards(InkyGloom),
					InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(3)))),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
					),
				},
			})

			h.P1.Play(InkyGloom)
			h.P1.EndTurn()
			h.P2.ChooseHouse(card.House.Brobnar)

			h.P2.ExpectCannotUseTo(foe, engine.ReapUse)
			h.P2.Fight(foe, ally) // fighting is still allowed
		},
	)

	t.Run("lifts once that turn is over", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(InkyGloom)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(InkyGloom)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Shadows)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P2.Reap(foe)
	})
}
