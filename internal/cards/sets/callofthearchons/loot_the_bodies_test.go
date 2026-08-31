package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Loot the Bodies
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time an enemy creature is destroyed, gain 1 Æmber.
func TestLootTheBodies(t *testing.T) {
	t.Run("gains 1 Æmber each time an enemy creature is destroyed this turn", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(LootTheBodies),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))))},
		})

		h.P1.Play(LootTheBodies)
		h.P1.Fight(ally, foe) // destroys the enemy creature

		h.Expect(foe).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
