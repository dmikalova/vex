package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Mars Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Mars card this turn.
func TestMarsAmbassador(t *testing.T) {
	t.Run("reap lets you play and use Mars cards this turn", func(t *testing.T) {
		var marsInPlay, marsInHand ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					MarsAmbassador,
					ct.Bind(&marsInPlay, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
				Hand: ct.Cards(
					ct.Bind(&marsInHand, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Reap(MarsAmbassador)

		if h.Game().State.MayPlayHouse[0] != engine.Mars {
			t.Error("Mars Ambassador should grant playing Mars cards this turn")
		}
		if h.Game().State.MayUseHouse[0] != engine.Mars {
			t.Error("Mars Ambassador should grant using friendly Mars creatures this turn")
		}

		// The play grant lets an out-of-house Mars card come down.
		h.P1.Play(marsInHand)
		h.Expect(marsInHand).At(ct.PlayArea)

		// The use grant lets an out-of-house Mars creature reap.
		h.P1.Reap(marsInPlay)
	})
}
