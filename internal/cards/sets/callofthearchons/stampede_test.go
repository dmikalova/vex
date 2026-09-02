package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Stampede
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If you used 3 or more creatures this turn, steal 2 Æmber.
func TestStampede(t *testing.T) {
	setup := func() ct.Setup {
		return ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Stampede),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Untamed)),
					ct.Creature(ct.OfHouse(card.House.Untamed)),
					ct.Creature(ct.OfHouse(card.House.Untamed)),
				),
			},
			P2: ct.Side{Amber: 5},
		}
	}

	t.Run("steals after three creatures were used", func(t *testing.T) {
		h := ct.Play(t, setup())
		for _, c := range h.Game().Battleline(0) {
			h.Game().ReapWith(c)
		}

		h.P1.Play(Stampede)

		h.P1.ExpectAmber(6) // 3 reaps, the Æmber bonus, and 2 stolen
		h.P2.ExpectAmber(3)
	})

	t.Run("steals nothing after only two creatures were used", func(t *testing.T) {
		h := ct.Play(t, setup())
		for _, c := range h.Game().Battleline(0)[:2] {
			h.Game().ReapWith(c)
		}

		h.P1.Play(Stampede)

		h.P1.ExpectAmber(3) // 2 reaps and the Æmber bonus
		h.P2.ExpectAmber(5)
	})
}
