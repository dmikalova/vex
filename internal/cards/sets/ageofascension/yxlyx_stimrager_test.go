package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yxlyx Stimrager
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Martian • Soldier
//
//	Fight: Deal 2 damage to a creature and move it to either flank of its controller's battleline.
func TestYxlyxStimrager(t *testing.T) {
	t.Run("its Fight ability damages a creature and moves it to a flank", func(t *testing.T) {
		var weak, a, b, c ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(YxlyxStimrager),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&weak, ct.Creature(ct.Power(1))),
					ct.Bind(&a, ct.Creature(ct.Power(6))),
					ct.Bind(&b, ct.Creature(ct.Power(6))),
					ct.Bind(&c, ct.Creature(ct.Power(6))),
				),
			},
		})

		h.P1.Fight(YxlyxStimrager, weak)
		h.P1.ClickCard(b)         // deal 2 damage to a surviving enemy
		h.P1.ClickOption("right") // move it to the right flank

		h.Expect(b).Damage(2)
		h.Expect(weak).At(ct.Discard)
		line := h.Game().Battleline(1)
		if len(line) != 3 || line[0] != a.ID() || line[1] != c.ID() || line[2] != b.ID() {
			t.Fatalf("battleline = %v, want [a c b]", line)
		}
	})
}
