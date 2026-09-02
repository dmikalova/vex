package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dr. Escotera
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Play: For each forged key your opponent has, gain 1 Æmber.
func TestDrEscotera(t *testing.T) {
	t.Run("gains 1 Æmber for each key the opponent has forged", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(DrEscotera)},
			P2: ct.Side{Keys: 2},
		})

		h.P1.Play(DrEscotera)

		h.P1.ExpectAmber(2)
	})
}
