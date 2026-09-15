package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Eater of the Dead
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Fight/Reap: Purge a creature from a discard pile -> give Eater of the Dead a +1 power counter.
func TestEaterOfTheDead(t *testing.T) {
	t.Run("reaping purges a creature from a discard pile and grows", func(t *testing.T) {
		var prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(EaterOfTheDead)},
			P2: ct.Side{Discard: ct.Cards(
				ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Brobnar))),
			)},
		})

		h.P1.Reap(EaterOfTheDead)

		h.Expect(prey).At(ct.Purge)
		h.Expect(EaterOfTheDead).Power(5) // 4 power + 1 counter
	})

	t.Run("does nothing with no creature to purge", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(EaterOfTheDead)},
		})

		h.P1.Reap(EaterOfTheDead)

		h.Expect(EaterOfTheDead).Power(4) // no counter placed
	})
}
