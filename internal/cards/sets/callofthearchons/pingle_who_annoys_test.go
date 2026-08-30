package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pingle Who Annoys
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin
//
//	Elusive.
//	Play: Deal 1 damage to each enemy creature.
func TestPingleWhoAnnoys(t *testing.T) {
	t.Run("deals 1 damage to each enemy creature when played", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(PingleWhoAnnoys)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Play(PingleWhoAnnoys)

		h.Expect(a).Damage(1)
		h.Expect(b).Damage(1)
	})
}
