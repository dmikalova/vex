package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tertiate
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy one third of all enemy creatures and one third of all friendly creatures, rounding up each time.
func TestTertiate(t *testing.T) {
	var f0, f1, f2 ct.Card
	var e0, e1, e2, e3 ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			InPlay: ct.Cards(
				ct.Bind(&f0, ct.Creature(ct.Power(3))),
				ct.Bind(&f1, ct.Creature(ct.Power(3))),
				ct.Bind(&f2, ct.Creature(ct.Power(3))),
			),
			Hand: ct.Cards(Tertiate),
		},
		P2: ct.Side{InPlay: ct.Cards(
			ct.Bind(&e0, ct.Creature(ct.Power(3))),
			ct.Bind(&e1, ct.Creature(ct.Power(3))),
			ct.Bind(&e2, ct.Creature(ct.Power(3))),
			ct.Bind(&e3, ct.Creature(ct.Power(3))),
		)},
	})

	h.P1.Play(Tertiate)
	// Enemy side: ceil(4/3)=2 destroyed; controller picks e0 and e2.
	h.P1.ClickCard(e0)
	h.P1.ClickCard(e2)
	// Friendly side: ceil(3/3)=1 destroyed; controller picks f0.
	h.P1.ClickCard(f0)

	h.Expect(e0).At(ct.Discard)
	h.Expect(e2).At(ct.Discard)
	h.Expect(e1).At(ct.PlayArea)
	h.Expect(e3).At(ct.PlayArea)
	h.Expect(f0).At(ct.Discard)
	h.Expect(f1).At(ct.PlayArea)
	h.Expect(f2).At(ct.PlayArea)
}
