package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gravelguts
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Giant
//
//	After a creature is destroyed fighting Gravelguts, give Gravelguts 2 +1 power counters.
func TestGravelguts(t *testing.T) {
	t.Run("gains two +1 power counters when an enemy is destroyed fighting it", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Gravelguts)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(2))))},
		})

		h.P1.Fight(Gravelguts, foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(Gravelguts).Power(7)
	})
}
