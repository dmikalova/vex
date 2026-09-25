package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Sinder
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Armor:  2
//	Traits: Demon
//
//	Taunt.
//	Reap: Destroy a friendly creature.
func TestSinder(t *testing.T) {
	t.Run("destroys a friendly creature when it reaps", func(t *testing.T) {
		var sinder, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&sinder, Sinder),
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.Reap(sinder)
		h.P1.ClickCard(ally)

		h.Expect(ally).At(ct.Discard)
	})
}
