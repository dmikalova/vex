package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yxilx Dominator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Armor:  1
//	Traits: Robot
//
//	Taunt.
//	Yxilx Dominator enters play stunned.
func TestYxilxDominator(t *testing.T) {
	t.Run("enters play stunned", func(t *testing.T) {
		var yxilx ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(ct.Bind(&yxilx, YxilxDominator))},
		})

		h.P1.Play(YxilxDominator)

		h.Expect(yxilx).At(ct.PlayArea).Stunned(true)
	})
}
