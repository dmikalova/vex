package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lupo the Scarred
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast
//
//	Skirmish.
//	Play: Deal 2 damage to an enemy creature.
func TestLupoTheScarred(t *testing.T) {
	t.Run("deals 2 damage to an enemy creature when played", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(LupoTheScarred)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))))},
		})

		h.P1.Play(LupoTheScarred)

		h.Expect(foe).At(ct.PlayArea).Damage(2)
	})
}
