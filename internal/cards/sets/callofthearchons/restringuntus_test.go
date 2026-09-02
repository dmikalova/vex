package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Restringuntus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: demon
//
//	Play: Choose a house - your opponent cannot choose that house as their active house until Restringuntus leaves play.
func TestRestringuntus(t *testing.T) {
	t.Run("bars the named house from the opponent until it leaves play", func(t *testing.T) {
		var restringuntus ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&restringuntus, Restringuntus)),
			},
			P2: ct.Side{House: card.House.Mars},
		})

		h.P1.Play(Restringuntus)
		h.P1.ExpectPrompt("Choose a house").Source("Restringuntus")
		h.P1.ClickOption("Mars")

		h.P1.EndTurn()
		if err := h.Game().ChooseHouse(1, card.House.Mars); err != engine.ErrHouseLocked {
			t.Errorf("choosing the barred house = %v, want ErrHouseLocked", err)
		}
		h.P2.ChooseHouse(card.House.Logos) // any other house is allowed
		h.P2.EndTurn()

		h.P1.ChooseHouse(card.House.Dis)
		h.Game().DestroyEach(0, []engine.LocalID{restringuntus.ID()})
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars) // the bar lifts with it
	})
}
