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
//	Play/Fight: You may purge a neighboring creature -> ready Buzzle.
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
		h.P1.ClickOption("Yes")

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
		h.P1.ClickOption("No")

		h.Expect(neighbor).At(ct.PlayArea)
		h.Expect(buzzle).Exhausted()
	})
}
