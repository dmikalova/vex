package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Brobnar Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Brobnar card this turn.
func TestBrobnarAmbassador(t *testing.T) {
	t.Run("reap grants playing and using Brobnar cards this turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(BrobnarAmbassador)},
		})

		h.P1.Reap(BrobnarAmbassador)

		if h.Game().State.MayPlayHouse[0] != engine.Brobnar {
			t.Error("Brobnar Ambassador should grant playing Brobnar cards this turn")
		}
		if h.Game().State.MayUseHouse[0] != engine.Brobnar {
			t.Error("Brobnar Ambassador should grant using friendly Brobnar creatures this turn")
		}
	})
}
