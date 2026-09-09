package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Information Exchange
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1A. If your opponent stole A from you on their previous turn,
//	steal 2A instead.
func TestInformationExchange(t *testing.T) {
	t.Run("opponent robbed you last turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(InformationExchange)},
			P2: ct.Side{House: card.House.Mars, Amber: 3},
		})
		h.Game().State.TurnHistory[0][engine.AemberStolenFromLastTurn] = 2

		h.P1.Play(InformationExchange)

		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(1)
	})

	t.Run("opponent did not rob you", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(InformationExchange)},
			P2: ct.Side{House: card.House.Mars, Amber: 3},
		})

		h.P1.Play(InformationExchange)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})
}
