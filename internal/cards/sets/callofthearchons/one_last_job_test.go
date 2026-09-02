package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// One Last Job
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Purge each friendly Shadows creature. For each creature purged this way, steal 1 Æmber.
func TestOneLastJob(t *testing.T) {
	t.Run("purges your Shadows creatures and steals for each", func(t *testing.T) {
		var shadows, brobnar, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(OneLastJob),
				InPlay: ct.Cards(
					ct.Bind(&shadows, ct.Creature(ct.OfHouse(card.House.Shadows))),
					ct.Bind(&brobnar, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
			P2: ct.Side{
				Amber:  5,
				InPlay: ct.Cards(ct.Bind(&theirs, ct.Creature(ct.OfHouse(card.House.Shadows)))),
			},
		})

		h.P1.Play(OneLastJob)

		h.Expect(shadows).At(ct.Purge)
		h.Expect(brobnar).At(ct.PlayArea) // only Shadows creatures are purged
		h.Expect(theirs).At(ct.PlayArea)  // only friendly creatures are purged
		h.P1.ExpectAmber(2)               // the Æmber bonus plus 1 stolen
		h.P2.ExpectAmber(4)
	})

	t.Run("steals nothing when you have no Shadows creatures", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(OneLastJob)},
			P2: ct.Side{Amber: 5},
		})

		h.P1.Play(OneLastJob)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(5)
	})
}
