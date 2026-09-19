package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Festering Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose up to 2 creatures. Deal 1 damage to each chosen creature. Deal 3 damage instead to each chosen creature that was already damaged.
func TestFesteringTouch(t *testing.T) {
	var clean, hurt ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(FesteringTouch),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&clean, ct.Creature(ct.Power(9))),
			ct.Bind(&hurt, ct.Creature(ct.Power(9))),
		)},
	})
	hurt.Damaged(1)

	h.P1.Play(FesteringTouch)
	h.P1.ClickCard(clean)
	h.P1.ClickCard(hurt)

	h.Expect(clean).Damage(1)
	h.Expect(hurt).Damage(4) // 1 already + 3 for being damaged
}
