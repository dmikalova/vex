package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Buzzle
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	Play/Fight: You may purge a neighboring Creature -> ready Buzzle.
func TestBuzzle(t *testing.T) {
	t.Run("purging a neighbor readies Buzzle when played", func(t *testing.T) {
		var buzzle, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&buzzle, Buzzle)),
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.Power(2))),
				),
			},
		})

		h.P1.Play(Buzzle)
		h.P1.ClickCard(neighbor)

		h.Expect(neighbor).At(ct.Purge)
		h.Expect(buzzle).Ready()
	})

	t.Run("declining leaves Buzzle exhausted", func(t *testing.T) {
		var buzzle, neighbor ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&buzzle, Buzzle)),
				InPlay: ct.Cards(
					ct.Bind(&neighbor, ct.Creature(ct.Power(2))),
				),
			},
		})

		h.P1.Play(Buzzle)
		h.P1.ClickDone()

		h.Expect(neighbor).At(ct.PlayArea)
		h.Expect(buzzle).Exhausted()
	})

	t.Run("no neighbor purges nothing without a prompt", func(t *testing.T) {
		var buzzle ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&buzzle, Buzzle)),
			},
		})

		h.P1.Play(Buzzle)

		h.Expect(buzzle).Exhausted()
	})
}
