package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hebe the Huge
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant • Knight
//
//	Play: Deal 2 damage to each other undamaged creature.
func TestHebeTheHuge(t *testing.T) {
	t.Run("deals 2 damage to each other undamaged creature", func(t *testing.T) {
		var healthy, hurt ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(HebeTheHuge)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&healthy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
					ct.Bind(&hurt, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		hurt.Damaged(1)
		h.P1.Play(HebeTheHuge)

		h.Expect(healthy).Damage(2)
		h.Expect(hurt).Damage(1)        // already damaged, so untouched
		h.Expect(HebeTheHuge).Damage(0) // itself excluded
	})
}
