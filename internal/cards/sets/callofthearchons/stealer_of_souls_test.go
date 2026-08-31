package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Stealer of Souls
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Demon
//
//	After a creature is destroyed fighting Stealer of Souls, purge it, and gain 1 Æmber.
func TestStealerOfSouls(t *testing.T) {
	t.Run("purges an enemy destroyed fighting it and gains 1 Æmber", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(StealerOfSouls)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))))},
		})

		h.P1.Fight(StealerOfSouls, foe)

		h.Expect(foe).At(ct.Purge)
		h.P1.ExpectAmber(1)
	})
}
