package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Shadows Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Shadows card this turn.
func TestShadowsAmbassador(t *testing.T) {
	t.Run("reap grants playing and using Shadows cards this turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(ShadowsAmbassador)},
		})

		h.P1.Reap(ShadowsAmbassador)

		if h.Game().State.MayPlayHouse[0] != engine.Shadows {
			t.Error("Shadows Ambassador should grant playing Shadows cards this turn")
		}
		if h.Game().State.MayUseHouse[0] != engine.Shadows {
			t.Error("Shadows Ambassador should grant using friendly Shadows creatures this turn")
		}
	})
}
