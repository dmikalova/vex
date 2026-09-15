package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hysteria
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Put each creature into its owner's hand.
func TestHysteria(t *testing.T) {
	t.Run("puts each creature into its owner's hand", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				Hand:   ct.Cards(Hysteria),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Dis)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Play(Hysteria)

		h.Expect(ally).At(ct.Hand)
		h.Expect(foe).At(ct.Hand)
	})
}
