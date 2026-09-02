package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Psychic Network
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each friendly ready Mars creature, steal 1 Æmber.
func TestPsychicNetwork(t *testing.T) {
	t.Run("steals 1 Æmber for each friendly ready Mars creature", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(PsychicNetwork),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Brobnar)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(PsychicNetwork)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(3)
	})

	t.Run("does not count an exhausted Mars creature", func(t *testing.T) {
		var spent ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(PsychicNetwork),
				InPlay: ct.Cards(
					ct.Bind(&spent, ct.Creature(ct.OfHouse(card.House.Mars))),
					ct.Creature(ct.OfHouse(card.House.Mars)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Reap(spent)
		h.P1.Play(PsychicNetwork)

		h.P1.ExpectAmber(2) // 1 reaped, 1 stolen for the one still ready
		h.P2.ExpectAmber(4)
	})
}
