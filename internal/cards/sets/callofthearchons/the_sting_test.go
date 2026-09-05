package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Sting
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Vehicle
//
//	You skip your "forge a key" step.
//	You gain all Æmber your opponent spends when forging a key.
//	Action: Destroy The Sting.
func TestTheSting(t *testing.T) {
	t.Run("gains the Æmber its controller's opponent spends forging a key", func(t *testing.T) {
		var sting ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{Amber: 6},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&sting, TheSting))},
		})

		h.P1.EndTurn() // P1's turn 1 ends; P2's turn 1 begins
		h.P2.EndTurn() // P2's turn ends; P1's turn 2 begins and P1 forges with its 6 Æmber

		h.P1.ExpectAmber(0)
		h.P1.ExpectKeys(1)
		h.P2.ExpectAmber(6)
	})

	t.Run("skips its controller's own forge step", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				Amber:  6,
				InPlay: ct.Cards(TheSting),
			},
		})

		h.P1.EndTurn() // P1's turn 1 ends; P2's turn 1 begins
		h.P2.EndTurn() // P2's turn ends; P1's turn 2 begins but skips forging

		h.P1.ExpectAmber(6)
		h.P1.ExpectKeys(0)
	})

	t.Run("can be destroyed", func(t *testing.T) {
		var sting ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&sting, TheSting)),
			},
		})

		h.P1.UseAction(sting)

		h.Expect(sting).At(ct.Discard)
	})
}
