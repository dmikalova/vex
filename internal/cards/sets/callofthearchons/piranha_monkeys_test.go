package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Piranha Monkeys
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast
//
//	Play/Reap: Deal 2 damage to each other creature.
func TestPiranhaMonkeys(t *testing.T) {
	t.Run("deals 2 damage to each other creature when played", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(PiranhaMonkeys)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Play(PiranhaMonkeys)

		h.Expect(a).Damage(2)
		h.Expect(b).Damage(2)
		h.Expect(PiranhaMonkeys).Damage(0) // itself excluded
	})
}
