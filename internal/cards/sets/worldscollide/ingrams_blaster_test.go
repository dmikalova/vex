package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ingram's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Medic Ingram, and fully heal a creature."
func TestIngramsBlaster(t *testing.T) {
	t.Run("attaches to Ingram and fully heals a creature", func(t *testing.T) {
		var carrier, wounded ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						IngramsBlaster,
					),
					MedicIngram,
					ct.Bind(
						&wounded,
						ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(6)),
					),
				),
			},
		})
		wounded.Damaged(4)

		h.P1.Reap(carrier)
		h.P1.ClickOption("Yes")
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")
		h.P1.ClickCard(wounded)

		h.Expect(wounded).Damage(0)
	})
}
