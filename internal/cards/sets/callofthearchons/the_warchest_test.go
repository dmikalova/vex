package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Warchest
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: For each enemy creature that was destroyed in a fight this turn, gain 1 Æmber.
func TestTheWarchest(t *testing.T) {
	t.Run("pays for each enemy creature killed in a fight this turn", func(t *testing.T) {
		var chest, brute, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&chest, TheWarchest),
					ct.Bind(&brute, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(9))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(1)))),
			},
		})

		h.P1.Fight(brute, enemy)
		h.Expect(enemy).At(ct.Discard)
		h.P1.UseAction(chest)
		h.P1.ExpectAmber(1)
	})

	t.Run("pays nothing when no enemy creature died fighting", func(t *testing.T) {
		var chest ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&chest, TheWarchest)),
			},
		})

		h.P1.UseAction(chest)
		h.P1.ExpectAmber(0)
	})
}
