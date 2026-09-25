package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mini Groupthink Tank
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Armor:  2
//	Traits: Robot • Experiment
//
//	Play/Fight/Reap: Deal 8 damage to a creature that shares a house with 2 of its neighbors.
func TestMiniGroupthinkTank(t *testing.T) {
	var a, b, c ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Logos,
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10))),
				ct.Bind(&c, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(10))),
			),
			Hand: ct.Cards(MiniGroupthinkTank),
		},
	})

	// Playing the tank enters it at a flank, leaving battleline A B C Tank. Only b
	// (middle Brobnar) shares a house with both of its neighbors, so it is the only
	// legal target and takes 8 damage automatically.
	h.P1.Play(MiniGroupthinkTank)

	h.Expect(b).Damage(8)
	h.Expect(a).Damage(0)
	h.Expect(c).Damage(0)
}
