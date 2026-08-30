package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Snufflegator
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast
//
//	Skirmish.
func TestSnufflegator(t *testing.T) {
	t.Run("takes no damage in return when it fights (skirmish)", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(Snufflegator)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)))),
			},
		})

		h.P1.Fight(Snufflegator, foe)

		h.Expect(Snufflegator).At(ct.PlayArea).Damage(0)
		h.Expect(foe).At(ct.PlayArea).Damage(4)
	})
}
