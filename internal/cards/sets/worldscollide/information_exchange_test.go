package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Information Exchange
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber. If your opponent stole Æmber from you on their previous turn, steal 1 Æmber.
func TestInformationExchange(t *testing.T) {
	t.Run("opponent robbed you last turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(InformationExchange),
			},
			P2: ct.Side{
				House: card.House.Mars,
				Amber: 3,
			},
		})
		h.Game().State.TurnHistory[0][engine.AemberStolenFromLastTurn] = 2

		h.P1.Play(InformationExchange)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("opponent did not rob you", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(InformationExchange),
			},
			P2: ct.Side{
				House: card.House.Mars,
				Amber: 3,
			},
		})

		h.P1.Play(InformationExchange)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})
}
