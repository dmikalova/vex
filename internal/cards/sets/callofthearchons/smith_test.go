package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Smith
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If you control more creatures than your opponent, gain 2 Æmber.
func TestSmith(t *testing.T) {
	t.Run("gains 2 Æmber when you control more creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(Smith),
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar))),
			},
		})

		h.P1.Play(Smith)

		h.P1.ExpectAmber(3) // 1 bonus + 2 from the condition
	})

	t.Run("only the bonus when creatures are even", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(Smith),
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Brobnar))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Mars)))},
		})

		h.P1.Play(Smith)

		h.P1.ExpectAmber(1) // bonus only
	})
}
