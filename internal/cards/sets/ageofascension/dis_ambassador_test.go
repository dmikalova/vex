package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Dis Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Dis card this turn.
func TestDisAmbassador(t *testing.T) {
	t.Run("reap grants playing and using Dis cards this turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(DisAmbassador)},
		})

		h.P1.Reap(DisAmbassador)

		if h.Game().State.MayPlayHouse[0] != engine.Dis {
			t.Error("Dis Ambassador should grant playing Dis cards this turn")
		}
		if h.Game().State.MayUseHouse[0] != engine.Dis {
			t.Error("Dis Ambassador should grant using friendly Dis creatures this turn")
		}
	})
}
