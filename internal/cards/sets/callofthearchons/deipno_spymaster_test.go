package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Deipno Spymaster
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive. Versatile.
//	Action: Use a friendly creature.
func TestDeipnoSpymaster(t *testing.T) {
	t.Run("Action uses a chosen friendly creature (Worker reaps)", func(t *testing.T) {
		var worker ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					DeipnoSpymaster,
					ct.Bind(&worker, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3))),
					// A second ready creature, so which creature is used is a real choice:
					// only ready creatures are offered, and Deipno exhausts to use its action.
					ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(3)),
				),
			},
		})

		h.P1.UseAction(DeipnoSpymaster)
		h.P1.ClickCard(worker) // use Worker, not the other creature

		h.P1.ExpectAmber(1) // Worker reaped for 1 Æmber
	})
}
