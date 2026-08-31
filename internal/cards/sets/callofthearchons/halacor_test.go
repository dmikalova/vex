package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Halacor
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Beast
//
//	Each friendly flank creature gains skirmish.
func TestHalacor(t *testing.T) {
	t.Run("grants skirmish to friendly flank creatures", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					Halacor,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(ally, foe)

		h.Expect(ally).At(ct.PlayArea).Damage(0) // skirmish granted to the flank creature
	})
}
