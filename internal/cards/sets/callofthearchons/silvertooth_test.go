package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Silvertooth
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Silvertooth enters play ready.
func TestSilvertooth(t *testing.T) {
	t.Run("enters play ready", func(t *testing.T) {
		var silver ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(ct.Bind(&silver, Silvertooth))},
		})

		h.P1.Play(Silvertooth)

		if silver.Exhausted() {
			t.Error("Silvertooth should enter play ready, not exhausted")
		}
	})
}
