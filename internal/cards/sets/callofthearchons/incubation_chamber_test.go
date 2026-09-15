package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Incubation Chamber
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Versatile.
//	Action: Reveal a Mars creature from your hand and archive it.
func TestIncubationChamber(t *testing.T) {
	t.Run("archives a Mars creature from hand", func(t *testing.T) {
		var chamber, martian, tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&chamber, IncubationChamber)),
				Hand: ct.Cards(
					ct.Bind(&tactic, ct.Tactic(ct.OfHouse(card.House.Mars))),
					ct.Bind(&martian, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(chamber)

		h.Expect(martian).At(ct.Archives)
		h.Expect(tactic).At(ct.Hand)
	})

	t.Run("does nothing with no Mars creature in hand", func(t *testing.T) {
		var chamber, offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&chamber, IncubationChamber)),
				Hand: ct.Cards(
					ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.UseAction(chamber)

		h.Expect(offHouse).At(ct.Hand)
	})
}
