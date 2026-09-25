package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Fission Bloom
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Power
//
//	Action: The next time you play a card this turn, resolve each of its bonus icons an additional time.
//	Enhance Draw.
func TestFissionBloom(t *testing.T) {
	t.Run("next played card resolves each bonus icon an additional time", func(t *testing.T) {
		var bloom, budded ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&bloom, FissionBloom)),
				Hand: ct.Cards(
					ct.Bind(&budded, ct.Creature(
						ct.OfHouse(card.House.Logos), ct.Bonus(card.Bonus.Aember))),
				),
			},
		})

		h.P1.UseAction(bloom)
		before := h.P1.Amber()
		h.P1.Play(budded)

		if got := h.P1.Amber() - before; got != 2 {
			t.Fatalf("aember gained = %d, want 2", got)
		}
	})
}
