package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chieftain's Brew
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Give a creature two +1 power counters.
func TestChieftainsBrew(t *testing.T) {
	t.Run("gives a creature two +1 power counters", func(t *testing.T) {
		var troll ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(ChieftainsBrew),
				InPlay: ct.Cards(ct.Bind(&troll, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Play(ChieftainsBrew)

		h.Expect(troll).Power(6)
	})
}
