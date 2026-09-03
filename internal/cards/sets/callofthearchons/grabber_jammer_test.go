package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Grabber Jammer
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Robot
//
//	Your opponent's keys cost +1 Æmber.
//	Fight/Reap: Grabber Jammer captures 1 Æmber from your opponent.
func TestGrabberJammer(t *testing.T) {
	t.Run("captures 1 Æmber when it fights or reaps", func(t *testing.T) {
		var jammer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(ct.Bind(&jammer, GrabberJammer))},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(jammer)

		h.Expect(jammer).AmberOn(1)
		h.P2.ExpectAmber(2) // 1 captured
	})

	t.Run("raises the opponent's key cost by 1 while in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(GrabberJammer)},
		})

		g := h.Game()
		g.State.Aember[1] = engine.KeyCost
		g.StartTurn(1)
		h.P2.ExpectKeys(0) // one Æmber short of the raised cost

		g.State.Aember[1] = engine.KeyCost + 1
		g.StartTurn(1)
		h.P2.ExpectKeys(1)
		h.P2.ExpectAmber(0) // paid the raised cost
	})
}
