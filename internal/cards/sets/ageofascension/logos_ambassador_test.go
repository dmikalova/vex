package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Logos Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Logos card this turn.
func TestLogosAmbassador(t *testing.T) {
	t.Run("reap grants playing and using Logos cards this turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(LogosAmbassador)},
		})

		h.P1.Reap(LogosAmbassador)

		if h.Game().State.MayPlayHouse[0] != engine.Logos {
			t.Error("Logos Ambassador should grant playing Logos cards this turn")
		}
		if h.Game().State.MayUseHouse[0] != engine.Logos {
			t.Error("Logos Ambassador should grant using friendly Logos creatures this turn")
		}
	})
}
