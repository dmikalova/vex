package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bull-wark
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Mutant • Knight
//
//	Assault 2.
//	Each neighboring creature gains assault 2.
func TestBullwark(t *testing.T) {
	var neighbor, foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			InPlay: ct.Cards(
				ct.Bind(&neighbor, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(1))),
				Bullwark,
			),
		},
		P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
	})

	// The neighbor gains assault 2 from Bull-wark, so fighting deals 2 assault
	// damage before combat, plus its own 1 power.
	h.P1.Fight(neighbor, foe)
	h.Expect(foe).Damage(3)
}
