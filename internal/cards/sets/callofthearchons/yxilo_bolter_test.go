package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Yxilo Bolter
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Fight/Reap: Deal 2 damage to a creature. If this damage destroys that creature, purge it.
func TestYxiloBolter(t *testing.T) {
	t.Run("purges a creature its damage destroys when it reaps", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(YxiloBolter)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(1))),
				),
			},
		})

		h.P1.Reap(YxiloBolter)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.Purge)
	})
}
