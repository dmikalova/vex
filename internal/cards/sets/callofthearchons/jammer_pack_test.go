package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Jammer Pack
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Your opponent's keys cost +2 Æmber."
func TestJammerPack(t *testing.T) {
	t.Run("raises the opponent's key cost by 2 while attached", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				Hand:   ct.Cards(JammerPack),
			},
		})

		h.P1.Play(JammerPack) // the lone host auto-attaches

		g := h.Game()
		// Seven Æmber is a key short: the raised cost is eight.
		g.State.Aember[1] = engine.KeyCost + 1
		g.StartTurn(1)
		h.P2.ExpectKeys(0)

		g.State.Aember[1] = engine.KeyCost + 2
		g.StartTurn(1)
		h.P2.ExpectKeys(1)
		h.P2.ExpectAmber(0) // paid the raised cost
	})
}
