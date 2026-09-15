package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Iron Obelisk
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Your opponent's keys cost +1 Æmber for each friendly damaged Brobnar creature.
func TestIronObelisk(t *testing.T) {
	t.Run("charges the opponent 1 more per friendly damaged Brobnar creature", func(t *testing.T) {
		var hurt, healthy, offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					IronObelisk,
					ct.Bind(&hurt, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(9))),
					ct.Bind(&healthy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(9))),
					ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(9))),
				),
			},
		})

		if got := h.Game().CurrentKeyCost(1); got != 6 {
			t.Errorf("key cost with no damaged Brobnar creature = %d, want 6", got)
		}

		hurt.Damaged(1)
		offHouse.Damaged(1) // wrong house, does not count
		if got := h.Game().CurrentKeyCost(1); got != 7 {
			t.Errorf("key cost with one damaged Brobnar creature = %d, want 7", got)
		}

		healthy.Damaged(1)
		if got := h.Game().CurrentKeyCost(1); got != 8 {
			t.Errorf("key cost with two damaged Brobnar creatures = %d, want 8", got)
		}
		if got := h.Game().CurrentKeyCost(0); got != 6 {
			t.Errorf("the controller's own key cost = %d, want 6", got)
		}
	})
}
