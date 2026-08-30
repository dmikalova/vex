package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Flaxia
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Faerie
//
//	Play: If you control more creatures than your opponent, gain 2 Æmber.
func TestFlaxia(t *testing.T) {
	t.Run("gains 2 Æmber when you control more creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(Flaxia)},
		})

		h.P1.Play(Flaxia) // Flaxia itself makes it 1 vs 0

		h.P1.ExpectAmber(2)
	})
}
