package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Groupthink Tank
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  3
//	Traits: Robot • Experiment
//
//	Action: Deal 4D to each creature that shares a house with at least 1 of its
//	neighbors.
func TestGroupthinkTank(t *testing.T) {
	var a, b, c ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
				ct.Bind(&c, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(6))),
				GroupthinkTank,
			),
		},
	})

	h.P1.UseAction(GroupthinkTank)

	// a and b are Brobnar neighbors, so each shares a house with a neighbor.
	h.Expect(a).Damage(4)
	h.Expect(b).Damage(4)
	// c (Shadows) neighbors a Brobnar and the Logos tank — no shared house.
	h.Expect(c).Damage(0)
	// The tank's only neighbor is Shadows, so it takes no damage.
	h.Expect(GroupthinkTank).Damage(0)
}
