package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Unnatural Selection
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose 3 friendly Creatures and 3 enemy Creatures - destroy each other Creature.
func TestUnnaturalSelection(t *testing.T) {
	var f0, f1, f2, f3 ct.Card
	var e0, e1, e2, e3 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Untamed,
			InPlay: ct.Cards(
				ct.Bind(&f0, ct.Creature(ct.Power(3))),
				ct.Bind(&f1, ct.Creature(ct.Power(3))),
				ct.Bind(&f2, ct.Creature(ct.Power(3))),
				ct.Bind(&f3, ct.Creature(ct.Power(3))),
			),
			Hand: ct.Cards(UnnaturalSelection),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&e0, ct.Creature(ct.Power(3))),
			ct.Bind(&e1, ct.Creature(ct.Power(3))),
			ct.Bind(&e2, ct.Creature(ct.Power(3))),
			ct.Bind(&e3, ct.Creature(ct.Power(3))),
		)},
	})

	h.P1.Play(UnnaturalSelection)
	// Keep f0,f1,f2 on the friendly side and e0,e1,e2 on the enemy side.
	h.P1.ClickCard(f0)
	h.P1.ClickCard(f1)
	h.P1.ClickCard(f2)
	h.P1.ClickCard(e0)
	h.P1.ClickCard(e1)
	h.P1.ClickCard(e2)

	h.Expect(f0).At(ct.PlayArea)
	h.Expect(f1).At(ct.PlayArea)
	h.Expect(f2).At(ct.PlayArea)
	h.Expect(f3).At(ct.Discard)
	h.Expect(e0).At(ct.PlayArea)
	h.Expect(e1).At(ct.PlayArea)
	h.Expect(e2).At(ct.PlayArea)
	h.Expect(e3).At(ct.Discard)
}
