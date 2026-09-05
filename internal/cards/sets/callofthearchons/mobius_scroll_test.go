package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mobius Scroll
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Archive Mobius Scroll from play, and archive up to 2 cards from your hand.
func TestMobiusScroll(t *testing.T) {
	t.Run("archives itself and two chosen cards", func(t *testing.T) {
		var scroll, first, second, kept ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&scroll, MobiusScroll)),
				Hand: ct.Cards(
					ct.Bind(&first, ct.Creature()),
					ct.Bind(&second, ct.Creature()),
					ct.Bind(&kept, ct.Creature()),
				),
			},
		})
		scroll.Ready()

		h.P1.UseAction(scroll)
		h.P1.ClickCard(first)
		h.P1.ClickCard(second)

		h.Expect(scroll).At(ct.Archives)
		h.Expect(first).At(ct.Archives)
		h.Expect(second).At(ct.Archives)
		h.Expect(kept).At(ct.Hand)
	})

	t.Run("archives fewer cards when the controller stops", func(t *testing.T) {
		var scroll, first, kept ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&scroll, MobiusScroll)),
				Hand: ct.Cards(
					ct.Bind(&first, ct.Creature()),
					ct.Bind(&kept, ct.Creature()),
				),
			},
		})
		scroll.Ready()

		h.P1.UseAction(scroll)
		h.P1.ClickCard(first)
		h.P1.ClickDone()

		h.Expect(first).At(ct.Archives)
		h.Expect(kept).At(ct.Hand)
	})
}
