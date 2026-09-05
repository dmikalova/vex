package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ortannu the Chained
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	Reap: Put each Ortannu's Binding from your discard pile into your hand. For each card returned this way, deal 2 damage to a creature that is not on a flank and 2 damage to each of its neighbors.
func TestOrtannuTheChained(t *testing.T) {
	t.Run("returns each Binding and deals a splash hit per return", func(t *testing.T) {
		var mid, left, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Dis,
				InPlay:  ct.Cards(OrtannuTheChained),
				Discard: ct.Cards(OrtannusBinding, OrtannusBinding),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(10))),
					ct.Bind(&mid, ct.Creature(ct.Power(10))),
					ct.Bind(&right, ct.Creature(ct.Power(10))),
				),
			},
		})

		h.P1.Reap(OrtannuTheChained)

		// Two returns → two spread hits on the non-flank creature: 2+2 = 4 damage,
		// with 2 splash per hit to each neighbor.
		h.Expect(mid).Damage(4)
		h.Expect(left).Damage(4)
		h.Expect(right).Damage(4)
	})

	t.Run("does nothing with no Binding in the discard pile", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(OrtannuTheChained),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(10))),
				),
			},
		})

		h.P1.Reap(OrtannuTheChained)

		h.Expect(foe).Damage(0)
	})
}
