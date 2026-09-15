package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mighty Tiger
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Play: Deal 4 damage to an enemy creature.
func TestMightyTiger(t *testing.T) {
	t.Run("deals 4 damage to an enemy creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(MightyTiger)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
				),
			},
		})

		h.P1.Play(MightyTiger)

		h.Expect(foe).At(ct.PlayArea).Damage(4)
	})
}
