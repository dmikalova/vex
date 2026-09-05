package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mating Season
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Shuffle each Mars creature into its owner's deck. For each creature shuffled into their deck this way, each player gains 1 Æmber.
func TestMatingSeason(t *testing.T) {
	mars := ct.OfHouse(card.House.Mars)

	t.Run(
		"shuffles each Mars creature away and pays each player for their own",
		func(t *testing.T) {
			var mine, mineOther, theirs1, theirs2 ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Mars,
					Hand:  ct.Cards(MatingSeason),
					InPlay: ct.Cards(
						ct.Bind(&mine, ct.Creature(mars)),
						ct.Bind(&mineOther, ct.Creature(ct.OfHouse(card.House.Dis))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&theirs1, ct.Creature(mars)),
						ct.Bind(&theirs2, ct.Creature(mars)),
					),
				},
			})

			h.P1.Play(MatingSeason)
			h.Expect(mine).At(ct.Deck)
			h.Expect(theirs1).At(ct.Deck)
			h.Expect(theirs2).At(ct.Deck)
			h.Expect(mineOther).At(ct.PlayArea)

			// One of their own back in the deck plus the card's own 1 Æmber bonus;
			// the opponent collects for the two creatures they controlled.
			h.P1.ExpectAmber(2)
			h.P2.ExpectAmber(2)
		},
	)
}
