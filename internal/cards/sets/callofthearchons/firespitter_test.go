package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Firespitter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Giant
//
//	Before Fight: Deal 1 damage to each enemy creature.
func TestFirespitter(t *testing.T) {
	t.Run("deals 1 damage to each enemy creature before fighting", func(t *testing.T) {
		var weak, tough ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Firespitter)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(1))),
					ct.Bind(&tough, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10))),
				),
			},
		})

		h.P1.Fight(Firespitter, tough)

		h.Expect(weak).At(ct.Discard) // destroyed by the before-fight damage
		h.Expect(tough).Damage(6)     // 1 before-fight + 5 combat
	})
}
