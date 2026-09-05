package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Untamed Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Untamed card this turn.
func TestUntamedAmbassador(t *testing.T) {
	t.Run("reap grants playing and using Untamed cards this turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(UntamedAmbassador)},
		})

		h.P1.Reap(UntamedAmbassador)

		if h.Game().State.MayPlayHouse[0] != engine.Untamed {
			t.Error("Untamed Ambassador should grant playing Untamed cards this turn")
		}
		if h.Game().State.MayUseHouse[0] != engine.Untamed {
			t.Error("Untamed Ambassador should grant using friendly Untamed creatures this turn")
		}
	})
}
