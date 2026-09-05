package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Groke
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Fight: Your opponent loses 1 Æmber.
func TestGroke(t *testing.T) {
	t.Run("makes the opponent lose 1 Æmber when it fights", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Groke)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3)))),
				Amber:  3,
			},
		})

		h.P1.Fight(Groke, enemy)

		h.P2.ExpectAmber(2)
	})
}
