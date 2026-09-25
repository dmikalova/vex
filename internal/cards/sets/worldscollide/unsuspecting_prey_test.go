package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Unsuspecting Prey
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2 damage to up to 3 undamaged creatures.
func TestUnsuspectingPrey(t *testing.T) {
	var clean1, clean2, hurt ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			Hand:  ct.Cards(UnsuspectingPrey),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&clean1, ct.Creature(ct.Power(6))),
			ct.Bind(&clean2, ct.Creature(ct.Power(6))),
			ct.Bind(&hurt, ct.Creature(ct.Power(6))),
		)},
	})
	hurt.Damaged(1)

	h.P1.Play(UnsuspectingPrey)
	h.P1.ClickCard(clean1)
	h.P1.ClickCard(clean2)

	h.Expect(clean1).Damage(2)
	h.Expect(clean2).Damage(2)
	h.Expect(hurt).Damage(1)
}
