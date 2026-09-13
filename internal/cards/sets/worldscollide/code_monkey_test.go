package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Code Monkey
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: AI • Beast
//
//	Deploy.
//	Play: Archive each neighboring Creature from play. If those Creatures share a house, gain 2 Æmber.
func TestCodeMonkey(t *testing.T) {
	t.Run("archives both neighbors and gains 2 when they share a house", func(t *testing.T) {
		var left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
				Hand: ct.Cards(CodeMonkey),
			},
		})

		h.P1.Play(CodeMonkey)
		h.P1.ClickOption("Between") // deploy between the two neighbors

		h.Expect(left).At(ct.Archives)
		h.Expect(right).At(ct.Archives)
		h.P1.ExpectAmber(2)
	})

	t.Run("archives both neighbors but gains nothing when houses differ", func(t *testing.T) {
		var left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Shadows))),
				),
				Hand: ct.Cards(CodeMonkey),
			},
		})

		h.P1.Play(CodeMonkey)
		h.P1.ClickOption("Between")

		h.Expect(left).At(ct.Archives)
		h.Expect(right).At(ct.Archives)
		h.P1.ExpectAmber(0)
	})
}
