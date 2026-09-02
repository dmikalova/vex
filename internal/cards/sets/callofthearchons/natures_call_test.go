package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nature's Call
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put up to 3 creatures into their owners' hands.
func TestNaturesCall(t *testing.T) {
	t.Run("puts up to 3 creatures into their owners' hands", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(NaturesCall),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Play(NaturesCall)
		h.P1.ClickCard(ally)
		h.P1.ClickCard(foe)

		h.Expect(ally).At(ct.Hand)
		h.Expect(foe).At(ct.Hand)
	})
}
