package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Amphora Captura
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	When resolving a bonus icon, you may resolve it as a Capture bonus icon instead.
//	Enhance Æmber Æmber Damage Damage Draw Draw.
func TestAmphoraCaptura(t *testing.T) {
	t.Run("may resolve an Æmber bonus icon as Capture instead", func(t *testing.T) {
		var bearer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(AmphoraCaptura),
				Hand: ct.Cards(
					ct.Bind(
						&bearer,
						ct.Creature(ct.OfHouse(card.House.Saurian), ct.Bonus(card.Bonus.Aember)),
					),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(bearer)
		h.P1.ClickOption("Yes") // capture instead of gain

		if got := h.P1.Amber(); got != 0 {
			t.Fatalf("P1 Æmber = %d, want 0 (captured, not gained)", got)
		}
		if got := h.P2.Amber(); got != 2 {
			t.Fatalf("P2 Æmber = %d, want 2 (1 captured)", got)
		}
	})
}
