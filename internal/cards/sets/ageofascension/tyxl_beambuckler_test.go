package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tyxl Beambuckler
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Martian • Soldier
//
//	Play: Deal 2 damage to a creature and move it to either flank of its controller's battleline.
func TestTyxlBeambuckler(t *testing.T) {
	t.Run("deals 2 damage to a creature and moves it to a flank", func(t *testing.T) {
		var a, b, c ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(TyxlBeambuckler),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.Power(6))),
					ct.Bind(&b, ct.Creature(ct.Power(6))),
					ct.Bind(&c, ct.Creature(ct.Power(6))),
				),
			},
		})

		h.P1.Play(TyxlBeambuckler)
		h.P1.ClickCard(b)         // deal 2 damage to the middle enemy
		h.P1.ClickOption("right") // move it to the right flank

		h.Expect(b).Damage(2)
		line := h.Game().Battleline(1)
		if len(line) != 3 || line[0] != a.ID() || line[1] != c.ID() || line[2] != b.ID() {
			t.Fatalf("battleline = %v, want [a c b]", line)
		}
	})
}
