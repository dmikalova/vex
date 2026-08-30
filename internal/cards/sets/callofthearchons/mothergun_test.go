package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mothergun
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Reveal any number of Mars cards from your hand, and for each card revealed this way, deal 1 damage to a creature.
func TestMothergun(t *testing.T) {
	t.Run("deals damage equal to the number of Mars cards revealed", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(Mothergun),
				Hand: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))))},
		})

		h.P1.UseAction(Mothergun)

		h.Expect(foe).At(ct.PlayArea).Damage(2)
	})
}
