package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Auto-Legionary
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Robot • Ally
//
//	Versatile.
//	Action: Give Auto-Legionary five +1 power counters. Move it to a flank of your battleline as a creature.
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

	t.Run("used again as a creature repositions instead of duplicating", func(t *testing.T) {
		var auto ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&auto, AutoLegionary)),
			},
		})

		h.P1.UseAction(auto)
		h.P1.ClickOption("left")

		// A second use finds it already a creature in the battleline; it must
		// reposition, not insert a duplicate — a duplicate breaks card conservation.
		auto.Ready()
		h.P1.UseAction(auto)
		h.P1.ClickOption("right")

		if err := h.Game().InvariantError(); err != nil {
			t.Fatalf("card conservation broken after repeated use: %v", err)
		}
		// It kept its counters across both uses (five each) and stayed one creature.
		h.Expect(auto).Power(10).At(ct.PlayArea)
	})
}
