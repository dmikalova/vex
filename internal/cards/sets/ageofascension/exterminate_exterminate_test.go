package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Exterminate! Exterminate!
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each non-Mars Creature with power less than the number of friendly Mars Creatures you control.
func TestExterminateExterminate(t *testing.T) {
	t.Run("destroys non-Mars creatures weaker than the friendly Mars count", func(t *testing.T) {
		var enemyWeak, enemyEqual, enemyMars, friendlyWeak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ExterminateExterminate),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Bind(
						&friendlyWeak,
						ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(2)),
					),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&enemyWeak, ct.Creature(ct.Power(2))),
				ct.Bind(&enemyEqual, ct.Creature(ct.Power(3))),
				ct.Bind(&enemyMars, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
			)},
		})

		// Three friendly Mars creatures set the threshold at 3, so every non-Mars
		// creature with power 2 or less falls; power 3 and the Mars creature stay.
		h.P1.Play(ExterminateExterminate)

		h.Expect(enemyWeak).At(ct.Discard)
		h.Expect(friendlyWeak).At(ct.Discard)
		h.Expect(enemyEqual).At(ct.PlayArea)
		h.Expect(enemyMars).At(ct.PlayArea)
	})
}
