package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Briar Grubbling
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast • Insect
//
//	Hazardous 5.
func TestBriarGrubbling(t *testing.T) {
	t.Run("Hazardous 5 destroys the attacker before any fight damage", func(t *testing.T) {
		var attacker ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(BriarGrubbling)},
		})

		h.P1.Fight(attacker, BriarGrubbling)

		h.Expect(attacker).At(ct.Discard)                  // destroyed by Hazardous 5
		h.Expect(BriarGrubbling).Damage(0).At(ct.PlayArea) // no combat occurred
	})
}
