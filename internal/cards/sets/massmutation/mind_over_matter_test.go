package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mind Over Matter
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Archive each creature from play.
func TestMindOverMatter(t *testing.T) {
	t.Run("archives each creature into its owner's archives", func(t *testing.T) {
		var mine, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				Hand:   ct.Cards(MindOverMatter),
				InPlay: ct.Cards(ct.Bind(&mine, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&theirs, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.Play(MindOverMatter)

		h.Expect(mine).At(ct.Archives)
		h.Expect(theirs).At(ct.Archives)
	})
}
