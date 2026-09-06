package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pain Reaction
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 2 damage to an enemy creature. If this damage destroys that creature, deal 2 damage to each of that creature's neighbors.
func TestPainReaction(t *testing.T) {
	t.Run("destroying the creature damages its neighbors", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(PainReaction),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(20))),
				ct.Bind(&middle, ct.Creature(ct.Power(2))),
				ct.Bind(&right, ct.Creature(ct.Power(20))),
			)},
		})

		h.P1.Play(PainReaction)
		h.P1.ClickCard(middle)

		h.Expect(middle).At(ct.Discard)
		h.Expect(left).Damage(2)
		h.Expect(right).Damage(2)
	})

	t.Run("surviving creature spares its neighbors", func(t *testing.T) {
		var left, middle, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(PainReaction),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.Power(20))),
				ct.Bind(&middle, ct.Creature(ct.Power(10))),
				ct.Bind(&right, ct.Creature(ct.Power(20))),
			)},
		})

		h.P1.Play(PainReaction)
		h.P1.ClickCard(middle)

		h.Expect(middle).Damage(2)
		h.Expect(left).Damage(0)
		h.Expect(right).Damage(0)
	})
}
