package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mogghunter's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains, "Fight: Deal 2 damage to a flank Creature."
func TestMogghuntersBrew(t *testing.T) {
	t.Run("host deals 2 damage to a flank creature when it fights", func(t *testing.T) {
		var host, bruiser, squishy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
						MogghuntersBrew,
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&bruiser, ct.Creature(ct.Power(3), ct.Armor(8))),
				ct.Bind(&squishy, ct.Creature(ct.Power(5))),
			)},
		})

		h.P1.Fight(host, bruiser)
		h.P1.ClickCard(squishy)

		h.Expect(squishy).Damage(2)
	})
}
