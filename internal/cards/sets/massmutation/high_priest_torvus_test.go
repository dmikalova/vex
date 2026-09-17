package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// High Priest Torvus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Priest
//
//	Reap: You may exalt High Priest Torvus -> after you resolve your next tactic this turn, put it into your hand instead of your discard pile.
func TestHighPriestTorvus(t *testing.T) {
	t.Run("exalts and returns the next action to hand", func(t *testing.T) {
		var torvus, action ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&torvus, HighPriestTorvus)),
				Hand: ct.Cards(
					ct.Bind(&action, ct.Tactic(ct.OfHouse(card.House.Saurian))),
				),
			},
		})

		h.P1.Reap(torvus)
		h.P1.ClickCard(torvus) // accept the optional exalt
		h.Expect(torvus).AmberOn(1)

		h.P1.Play(action)
		h.Expect(action).At(ct.Hand)
	})

	t.Run("declining the exalt discards the next action normally", func(t *testing.T) {
		var torvus, action ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&torvus, HighPriestTorvus)),
				Hand: ct.Cards(
					ct.Bind(&action, ct.Tactic(ct.OfHouse(card.House.Saurian))),
				),
			},
		})

		h.P1.Reap(torvus)
		h.P1.ClickDone() // decline the exalt
		h.Expect(torvus).AmberOn(0)

		h.P1.Play(action)
		h.Expect(action).At(ct.Discard)
	})
}
