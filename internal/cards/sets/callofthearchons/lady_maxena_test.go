package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lady Maxena
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Knight • Spirit
//
//	Play: Stun a creature.
//	Action: Put Lady Maxena into its owner's hand.
func TestLadyMaxena(t *testing.T) {
	t.Run("stuns a chosen creature when played", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(LadyMaxena)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Play(LadyMaxena)
		h.P1.ClickCard(foe)

		h.Expect(foe).Stunned(true)
	})

	t.Run("returns itself to hand as an action", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(LadyMaxena)},
		})

		h.P1.UseAction(LadyMaxena)

		h.Expect(LadyMaxena).At(ct.Hand)
	})
}
