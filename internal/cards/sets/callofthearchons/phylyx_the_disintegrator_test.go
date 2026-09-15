package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Phylyx the Disintegrator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Martian • Soldier
//
//	Elusive.
//	Action: For each other friendly Mars creature, your opponent loses 1 Æmber.
func TestPhylyxTheDisintegrator(t *testing.T) {
	t.Run("drains 1 Æmber per other friendly Mars creature", func(t *testing.T) {
		var phylyx ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&phylyx, PhylyxTheDisintegrator),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3)),
					ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)),
				),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(phylyx)

		// Two other Mars creatures; the Logos creature and Phylyx itself do not count.
		h.P2.ExpectAmber(3)
	})

	t.Run("drains nothing on its own", func(t *testing.T) {
		var phylyx ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&phylyx, PhylyxTheDisintegrator)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.P1.UseAction(phylyx)

		h.P2.ExpectAmber(5)
	})
}
