package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Primus Unguis
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Each friendly Creature gains +2 power for each Æmber on Primus Unguis.
//	Reap: Exalt Primus Unguis.
func TestPrimusUnguis(t *testing.T) {
	t.Run("every friendly creature grows with the Æmber on Primus Unguis", func(t *testing.T) {
		var primus, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&primus, PrimusUnguis),
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(3))),
				),
			},
		})

		// No Æmber on Primus yet: neither creature is buffed.
		h.Expect(primus).Power(5)
		h.Expect(friend).Power(3)

		// Reap exalts Primus, placing one Æmber on it, so every friendly creature
		// gains +2 power.
		h.P1.Reap(primus)

		h.Expect(primus).AmberOn(1)
		h.Expect(primus).Power(7)
		h.Expect(friend).Power(5)
	})
}
