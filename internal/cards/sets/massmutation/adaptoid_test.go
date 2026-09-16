package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Adaptoid
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	After you play a card with a bonus icon, choose one:
//	- For the remainder of the turn, Adaptoid gains +2 armor
//	- For the remainder of the turn, Adaptoid gains assault 2
//	- For the remainder of the turn, Adaptoid gains, "Fight: Steal 1 Æmber."
//	Enhance Capture Damage Draw.
func TestAdaptoid(t *testing.T) {
	t.Run("+2 armor option", func(t *testing.T) {
		var iconed ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Adaptoid),
				Hand: ct.Cards(
					ct.Bind(&iconed, ct.Creature(
						ct.OfHouse(card.House.Logos),
						ct.Power(3),
						ct.Bonus(card.Bonus.Aember))),
				),
			},
		})

		h.P1.Play(iconed)
		h.P1.ClickOption("armor")

		h.Expect(Adaptoid).Armor(2)
	})

	t.Run("Fight: Steal option", func(t *testing.T) {
		var iconed, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Adaptoid),
				Hand: ct.Cards(
					ct.Bind(&iconed, ct.Creature(
						ct.OfHouse(card.House.Logos),
						ct.Power(3),
						ct.Bonus(card.Bonus.Draw))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(5)))),
				Amber:  3,
			},
		})

		h.P1.Play(iconed)
		h.P1.ClickOption("steal")
		h.P1.Fight(Adaptoid, foe)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("a card without a bonus icon offers no choice", func(t *testing.T) {
		var plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Adaptoid),
				Hand: ct.Cards(
					ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(plain)

		h.Expect(Adaptoid).Armor(0)
	})
}
