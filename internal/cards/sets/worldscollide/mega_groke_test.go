package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mega Groke
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Fight: Your opponent loses 1 Æmber.
func TestMegaGroke(t *testing.T) {
	t.Run("opponent loses 1 Æmber when it fights", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(MegaGroke),
			},
			P2: ct.Side{
				Amber:  3,
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Fight(MegaGroke, foe)

		h.P2.ExpectAmber(2)
	})
}
