package massmutation

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// "Borrow"
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Take control of an enemy artifact. It belongs to house Shadows.
func TestBorrow(t *testing.T) {
	t.Run("takes control of an enemy artifact and reassigns it to Shadows", func(t *testing.T) {
		var relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(Borrow)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(Borrow)

		if !slices.Contains(h.Game().Artifacts(0), relic.ID()) {
			t.Error("the artifact should be controlled by P1")
		}
		if got := h.Game().House(relic.ID()); got != card.House.Shadows {
			t.Errorf("artifact house = %v, want Shadows", got)
		}
	})
}
