package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Knoxx
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Knoxx gains +3 power for each neighbor it has.
func TestKnoxx(t *testing.T) {
	t.Run("gets +3 power for each neighbor", func(t *testing.T) {
		var knoxx ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Creature(ct.Power(2)),
					ct.Bind(&knoxx, Knoxx),
					ct.Creature(ct.Power(2)),
				),
			},
		})

		h.Expect(knoxx).Power(9) // 3 + 3 per neighbor, two neighbors
	})

	t.Run("gets no bonus with no neighbors", func(t *testing.T) {
		var knoxx ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&knoxx, Knoxx)),
			},
		})

		h.Expect(knoxx).Power(3)
	})
}
