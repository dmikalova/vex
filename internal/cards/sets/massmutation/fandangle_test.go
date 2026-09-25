package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Fandangle
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Witch
//
//	While you have 4 or more Æmber, your non-Untamed creatures enter play ready.
func TestFandangle(t *testing.T) {
	t.Run("a non-Untamed creature enters ready while you hold 4 Æmber", func(t *testing.T) {
		var newbie ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				Amber:  4,
				InPlay: ct.Cards(Fandangle),
				Hand: ct.Cards(
					ct.Bind(&newbie, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(newbie)

		h.Expect(newbie).Ready()
	})

	t.Run("below 4 Æmber the grant is dormant", func(t *testing.T) {
		var newbie ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				Amber:  3,
				InPlay: ct.Cards(Fandangle),
				Hand: ct.Cards(
					ct.Bind(&newbie, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(newbie)

		h.Expect(newbie).Exhausted()
	})

	t.Run("an Untamed creature is excluded even while rich", func(t *testing.T) {
		var newbie ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Amber:  4,
				InPlay: ct.Cards(Fandangle),
				Hand: ct.Cards(
					ct.Bind(&newbie, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Play(newbie)

		h.Expect(newbie).Exhausted()
	})
}
