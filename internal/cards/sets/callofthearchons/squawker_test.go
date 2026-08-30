package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Squawker
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose one:
//	- Ready a Mars creature
//	- Stun a non-Mars creature.
func TestSquawker(t *testing.T) {
	t.Run("can stun a non-Mars creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(Squawker)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.Play(Squawker)
		h.P1.ClickOption("Stun a non-Mars creature")

		h.Expect(foe).Stunned(true)
	})
}
