package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Throwing Stars
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 1 damage to up to 3 Creatures. For each Creature destroyed this way, gain 1 Æmber.
func TestThrowingStars(t *testing.T) {
	t.Run("gains 1 Æmber per creature it destroys", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(ThrowingStars)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.Power(1))),
				ct.Bind(&b, ct.Creature(ct.Power(1))),
			)},
		})

		h.P1.Play(ThrowingStars)
		h.P1.ClickCard(a)
		h.P1.ClickCard(b)
		h.P1.ExpectAmber(2)
	})
}
