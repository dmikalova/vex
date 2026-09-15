package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Alaka's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight: Play a creature -> ready it."
func TestAlakasBrew(t *testing.T) {
	t.Run("host fighting plays a creature and readies it", func(t *testing.T) {
		var host, recruit, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand: ct.Cards(
					ct.Bind(&recruit, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
						AlakasBrew,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(host, foe)

		// The granted Fight ability played the recruit from hand and readied it.
		h.Expect(recruit).Ready()
	})
}
