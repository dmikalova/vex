package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Soft Landing
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: The next creature or artifact you play this turn enters play ready.
func TestSoftLanding(t *testing.T) {
	mars := ct.OfHouse(card.House.Mars)

	t.Run("the next creature enters play ready", func(t *testing.T) {
		var hound ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(SoftLanding, ct.Bind(&hound, ct.Creature(mars))),
			},
		})

		h.P1.Play(SoftLanding)
		h.P1.Play(hound)
		h.Expect(hound).Ready()
	})

	t.Run("the next artifact enters play ready", func(t *testing.T) {
		var sigil ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(SoftLanding, ct.Bind(&sigil, ct.Artifact(mars))),
			},
		})

		h.P1.Play(SoftLanding)
		h.P1.Play(sigil)
		h.Expect(sigil).Ready()
	})

	t.Run("only the next card benefits", func(t *testing.T) {
		var first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand: ct.Cards(
					SoftLanding,
					ct.Bind(&first, ct.Creature(mars)),
					ct.Bind(&second, ct.Creature(mars)),
				),
			},
		})

		h.P1.Play(SoftLanding)
		h.P1.Play(first)
		h.P1.Play(second)
		h.Expect(first).Ready()
		h.Expect(second).Exhausted()
	})
}
