package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Proclamation 346E
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	While fewer than 3 houses are represented among enemy creatures, your opponent's keys cost +2 Æmber.
func TestProclamation346E(t *testing.T) {
	t.Run("taxes the opponent while they field fewer than three houses", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(Proclamation346E),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Sanctum)),
					ct.Creature(ct.OfHouse(card.House.Untamed)),
				),
			},
		})

		if got := h.Game().CurrentKeyCost(1); got != 6+2 {
			t.Errorf("opponent key cost = %d, want 8", got)
		}
	})

	t.Run("stops taxing once the opponent fields three houses", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(Proclamation346E),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Sanctum)),
					ct.Creature(ct.OfHouse(card.House.Untamed)),
					ct.Creature(ct.OfHouse(card.House.Logos)),
				),
			},
		})

		if got := h.Game().CurrentKeyCost(1); got != 6 {
			t.Errorf("opponent key cost = %d, want 6", got)
		}
	})
}
