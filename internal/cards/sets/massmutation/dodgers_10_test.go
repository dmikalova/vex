package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Dodger's 10
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  11
//	Traits: Elf • Thief
//
//	Play/Fight/Reap: Steal half of your opponent's Æmber, rounded down.
func TestDodgers10(t *testing.T) {
	var dodgers ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Shadows,
			Hand:  ct.Cards(ct.Bind(&dodgers, Dodgers10), card.GiganticArt(Dodgers10)),
		},
		P2: ct.Side{Amber: 5},
	})

	h.P1.Play(dodgers)

	// Half of the opponent's 5 Æmber, rounded down, is 2.
	h.P2.ExpectAmber(3)
	h.P1.ExpectAmber(2)
}
