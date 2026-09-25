package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Equalize
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Redistribute the Æmber on friendly creatures among friendly creatures. Redistribute the Æmber on enemy creatures among enemy creatures.
func TestEqualize(t *testing.T) {
	t.Run("moves friendly Æmber among friendly creatures", func(t *testing.T) {
		var a, b, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(Equalize),
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
			},
		})

		h.Game().AddAmberOn(a.ID(), 2)
		h.Game().AddAmberOn(foe.ID(), 2)

		h.P1.Play(Equalize)
		// Redistributing the 2 friendly Æmber: place both onto b.
		h.P1.ClickCard(b)
		h.P1.ClickCard(b)

		h.Expect(a).AmberOn(0)
		h.Expect(b).AmberOn(2)
		// The lone enemy creature keeps its Æmber.
		h.Expect(foe).AmberOn(2)
	})
}
