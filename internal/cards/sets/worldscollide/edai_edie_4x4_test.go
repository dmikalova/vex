package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// EDAI "Edie" 4x4
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: AI • Scientist
//
//	Your opponent's keys cost +1 Æmber for each card in your archives.
//	Play: Archive a card from your hand.
func TestEDAIEdie4x4(t *testing.T) {
	t.Run("Play archives a card from hand", func(t *testing.T) {
		var edie, spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand: ct.Cards(
					ct.Bind(&edie, EDAIEdie4x4),
					ct.Bind(&spare, ct.Creature()),
				),
			},
		})

		h.P1.Play(edie)

		h.Expect(spare).At(ct.Archives)
	})

	t.Run("raises the opponent's key cost per card in your archives", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:    card.House.Logos,
				InPlay:   ct.Cards(EDAIEdie4x4),
				Archives: ct.Cards(ct.Creature(), ct.Creature()),
			},
		})

		base := 6
		if got := h.Game().CurrentKeyCost(1); got != base+2 {
			t.Errorf("opponent key cost = %d, want %d", got, base+2)
		}
		if got := h.Game().CurrentKeyCost(0); got != base {
			t.Errorf("controller key cost = %d, want %d", got, base)
		}
	})
}
