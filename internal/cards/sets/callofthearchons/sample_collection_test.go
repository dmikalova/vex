package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sample Collection
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each forged key your opponent has, put an enemy creature into your archives.
func TestSampleCollection(t *testing.T) {
	t.Run("abducts one enemy creature per key the opponent has forged", func(t *testing.T) {
		var collection, first, second, third ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&collection, SampleCollection)),
			},
			P2: ct.Side{
				Keys: 2,
				InPlay: ct.Cards(
					ct.Bind(&first, ct.Creature(ct.Power(3))),
					ct.Bind(&second, ct.Creature(ct.Power(4))),
					ct.Bind(&third, ct.Creature(ct.Power(5))),
				),
			},
		})

		h.P1.Play(collection)
		h.P1.ClickCard(first)
		h.P1.ClickCard(second)

		h.Expect(first).At(ct.Archives)
		h.Expect(second).At(ct.Archives)
		h.Expect(third).At(ct.PlayArea)
	})

	t.Run("abducts nothing when the opponent has forged no keys", func(t *testing.T) {
		var collection, only ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&collection, SampleCollection)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&only, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(collection)
		h.Expect(only).At(ct.PlayArea)
	})
}
