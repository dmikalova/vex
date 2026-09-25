package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Legion's March
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, after you use a Dinosaur creature, deal 1 damage to each non-Dinosaur creature.
func TestLegionsMarch(t *testing.T) {
	var march, dino, plain ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(ct.Bind(&march, LegionsMarch)),
			InPlay: ct.Cards(
				ct.Bind(&dino, ct.Creature(
					ct.OfHouse(card.House.Saurian),
					ct.Traits(card.Traits.Dinosaur),
					ct.Power(5),
				)),
				ct.Bind(&plain, ct.Creature(
					ct.OfHouse(card.House.Saurian),
					ct.Power(5),
				)),
			),
		},
	})

	h.P1.Play(march)
	h.P1.Reap(dino)

	h.Expect(plain).Damage(1)
	h.Expect(dino).Damage(0)
}
