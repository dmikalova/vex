package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Harland Mindlock
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Cyborg • Scientist
//
//	Play: Take control of an enemy flank creature until Harland Mindlock leaves play.
func TestHarlandMindlock(t *testing.T) {
	t.Run("takes control of an enemy flank creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(HarlandMindlock)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))))},
		})

		h.P1.Play(HarlandMindlock)

		if got := h.Game().Owner(foe.ID()); got != 1 {
			t.Errorf("ownership should stay with player 1, got %d", got)
		}
		inMine := false
		for _, id := range h.Game().Battleline(0) {
			if id == foe.ID() {
				inMine = true
			}
		}
		if !inMine {
			t.Error("the seized creature should be in the controller's battleline")
		}
	})
}
