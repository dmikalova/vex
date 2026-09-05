package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Pitlord
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  9
//	Æmber:  2
//	Traits: Demon
//
//	Taunt.
//	While Pitlord is in play you must choose Dis as your active house.
func TestPitlord(t *testing.T) {
	t.Run("its controller must choose Dis while it is in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Pitlord)},
			P2: ct.Side{House: card.House.Mars},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars)
		h.P2.EndTurn()

		if err := h.Game().ChooseHouse(0, card.House.Brobnar); err != engine.ErrHouseLocked {
			t.Errorf("choosing another house = %v, want ErrHouseLocked", err)
		}
		h.P1.ChooseHouse(card.House.Dis) // the locked house is allowed
	})

	t.Run("it does not lock the opponent's choice", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Pitlord)},
			P2: ct.Side{House: card.House.Mars},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars)
	})
}
