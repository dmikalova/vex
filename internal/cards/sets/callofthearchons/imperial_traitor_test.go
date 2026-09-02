package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Imperial Traitor
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Reveal your opponent's hand, and you may purge a Sanctum card from your opponent's hand.
func TestImperialTraitor(t *testing.T) {
	t.Run("purges a chosen Sanctum card from the opponent's hand", func(t *testing.T) {
		var sanctum, shadows ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(ImperialTraitor)},
			P2: ct.Side{Hand: ct.Cards(
				ct.Bind(&sanctum, ct.Creature(ct.OfHouse(card.House.Sanctum))),
				ct.Bind(&shadows, ct.Creature(ct.OfHouse(card.House.Shadows))),
			)},
		})

		h.P1.Play(ImperialTraitor)
		h.P1.ClickCard(sanctum)

		h.Expect(sanctum).At(ct.Purge)
		h.Expect(shadows).At(ct.Hand) // non-Sanctum cards are not eligible
	})

	t.Run("does nothing when the opponent reveals no Sanctum card", func(t *testing.T) {
		var shadows ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(ImperialTraitor)},
			P2: ct.Side{Hand: ct.Cards(
				ct.Bind(&shadows, ct.Creature(ct.OfHouse(card.House.Shadows))),
			)},
		})

		h.P1.Play(ImperialTraitor)

		h.Expect(shadows).At(ct.Hand)
	})
}
