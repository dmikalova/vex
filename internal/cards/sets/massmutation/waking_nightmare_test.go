package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Waking Nightmare
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Keys cost +1 Æmber for each Dis creature in play during your opponent's next turn.
//	Enhance Capture.
func TestWakingNightmare(t *testing.T) {
	t.Run("taxes the opponent per Dis creature in play", func(t *testing.T) {
		var nightmare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&nightmare, WakingNightmare)),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Dis)),
					ct.Creature(ct.OfHouse(card.House.Dis)),
				),
			},
			P2: ct.Side{House: card.House.Dis},
		})

		h.P1.Play(nightmare)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)

		// Two Dis creatures in play -> keys cost 6 + 2 during the opponent's turn.
		if got := h.Game().CurrentKeyCost(1); got != 8 {
			t.Errorf("key cost = %d, want 8", got)
		}
	})

	t.Run("the tax lifts after that turn", func(t *testing.T) {
		var nightmare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&nightmare, WakingNightmare)),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Dis)),
					ct.Creature(ct.OfHouse(card.House.Dis)),
				),
			},
			P2: ct.Side{House: card.House.Dis},
		})

		h.P1.Play(nightmare)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Dis)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)

		if got := h.Game().CurrentKeyCost(1); got != 6 {
			t.Errorf("key cost = %d, want 6 after the tax lifts", got)
		}
	})
}
