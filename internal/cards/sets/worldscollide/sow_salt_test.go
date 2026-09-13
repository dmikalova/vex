package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Sow Salt
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Alpha.
//	Play: Until the start of your next turn, Creatures cannot be used to reap.
func TestSowSalt(t *testing.T) {
	t.Run("bars every creature from reaping until the caster's next turn", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(SowSalt),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		// Alpha: Sow Salt is the first thing the caster does this turn.
		h.P1.Play(SowSalt)
		// The caster's own creatures cannot reap this turn.
		h.P1.ExpectCannotUseTo(ally, engine.ReapUse)
		h.P1.EndTurn()

		// The opponent's creatures cannot reap on their turn either.
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.ExpectCannotUseTo(foe, engine.ReapUse)
	})

	t.Run("lifts once the caster's next turn begins", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Saurian, Hand: ct.Cards(SowSalt)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(SowSalt)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.ExpectCannotUseTo(foe, engine.ReapUse)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Saurian)
		h.P1.EndTurn()

		// A full round later the bar is gone and reaping works again.
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Reap(foe)
	})
}
