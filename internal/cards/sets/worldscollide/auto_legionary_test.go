package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Auto-Legionary
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Robot • Ally
//
//	Versatile.
//	Action: Give Auto-Legionary five +1 power counters. Move it to a flank of your battleline as a Creature.
func TestAutoLegionary(t *testing.T) {
	t.Run("turns itself into a 5-power creature on the chosen flank", func(t *testing.T) {
		var auto, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.Power(3))),
					ct.Bind(&auto, AutoLegionary),
				),
			},
		})

		// Before its action it is an artifact.
		if got := h.Game().TypeOf(auto.ID()); got != engine.Artifact {
			t.Fatalf("Auto-Legionary should start as an artifact, got %v", got)
		}

		h.P1.UseAction(auto)
		h.P1.ClickOption("right") // move it to the right flank

		if got := h.Game().TypeOf(auto.ID()); got != engine.Creature {
			t.Fatalf("Auto-Legionary should have become a creature, got %v", got)
		}
		h.Expect(auto).Power(5).Exhausted().At(ct.PlayArea)
		// It joins the battleline; the existing creature is unchanged.
		h.Expect(ally).Power(3).At(ct.PlayArea)
	})

	t.Run("Versatile lets it act out of house", func(t *testing.T) {
		var auto ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&auto, AutoLegionary)),
			},
		})

		h.P1.UseAction(auto)
		h.P1.ClickOption("left")

		if got := h.Game().TypeOf(auto.ID()); got != engine.Creature {
			t.Fatalf("Auto-Legionary should have become a creature, got %v", got)
		}
		h.Expect(auto).Power(5).At(ct.PlayArea)
	})
}
