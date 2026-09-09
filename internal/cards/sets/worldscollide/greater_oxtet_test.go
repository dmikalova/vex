package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Greater Oxtet
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Taunt.
//	At the end of your "ready cards" step, purge a card from your hand -> give Greater Oxtet two +1 power counters.
func TestGreaterOxtet(t *testing.T) {
	t.Run("purges a card to grow by two", func(t *testing.T) {
		var fodder ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(GreaterOxtet),
				Hand:   ct.Cards(ct.Bind(&fodder, ct.Tactic(ct.OfHouse(card.House.Dis)))),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()

		h.Expect(fodder).At(ct.Purge)
		h.Expect(GreaterOxtet).Power(6)
	})

	t.Run("an empty hand leaves it unchanged", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(GreaterOxtet),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()

		h.Expect(GreaterOxtet).Power(4)
	})
}
