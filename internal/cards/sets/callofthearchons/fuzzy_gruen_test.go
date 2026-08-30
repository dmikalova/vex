package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fuzzy Gruen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Æmber:  2
//	Traits: Beast
//
//	Play: Your opponent gains 1 Æmber.
func TestFuzzyGruen(t *testing.T) {
	t.Run("gives the opponent 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(FuzzyGruen)},
			P2: ct.Side{Amber: 0},
		})

		h.P1.Play(FuzzyGruen)

		h.P2.ExpectAmber(1)
	})
}
