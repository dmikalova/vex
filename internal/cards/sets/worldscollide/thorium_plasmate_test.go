package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Thorium Plasmate
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Move an enemy creature anywhere in its controller's battleline. Deal 2D to that creature for each of its neighbors that shares a house with it.
func TestThoriumPlasmate(t *testing.T) {
	t.Run("deals 2 per neighbor sharing the moved creature's house", func(t *testing.T) {
		var moved, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(ThoriumPlasmate)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&moved, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
				),
			},
		})

		h.P1.Play(ThoriumPlasmate)
		h.P1.ClickCard(moved)
		h.P1.ClickOption("Left flank")

		// The moved creature's only neighbor shares its house, so it takes 2 damage.
		h.Expect(moved).Damage(2)
	})

	t.Run("deals nothing when no neighbor shares the moved creature's house", func(t *testing.T) {
		var moved, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(ThoriumPlasmate)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&moved, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(6))),
				),
			},
		})

		h.P1.Play(ThoriumPlasmate)
		h.P1.ClickCard(moved)
		h.P1.ClickOption("Left flank")

		h.Expect(moved).Damage(0)
	})
}
