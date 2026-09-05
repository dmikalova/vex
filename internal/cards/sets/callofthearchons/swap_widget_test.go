package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Swap Widget
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Put a friendly ready Mars creature into its owner's hand -> put a Mars creature with a different name from your hand into play, and ready it.
func TestSwapWidget(t *testing.T) {
	t.Run("returns a ready Mars creature and swaps in a differently named one", func(t *testing.T) {
		var widget, tunk, sameName, blypyp ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&widget, SwapWidget),
					ct.Bind(&tunk, Tunk),
				),
				Hand: ct.Cards(
					ct.Bind(&sameName, Tunk),
					ct.Bind(&blypyp, Blypyp),
				),
			},
		})

		h.P1.UseAction(widget)

		h.Expect(tunk).At(ct.Hand)
		h.Expect(sameName).At(ct.Hand)
		h.Expect(blypyp).At(ct.PlayArea).Ready()
	})

	t.Run("with no ready friendly Mars creature, nothing swaps", func(t *testing.T) {
		var widget, exhausted, blypyp ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&widget, SwapWidget),
					ct.Bind(&exhausted, Tunk),
				),
				Hand: ct.Cards(ct.Bind(&blypyp, Blypyp)),
			},
		})
		h.Game().SetExhausted(exhausted.ID(), true)

		h.P1.UseAction(widget)

		h.Expect(exhausted).At(ct.PlayArea)
		h.Expect(blypyp).At(ct.Hand)
	})
}
