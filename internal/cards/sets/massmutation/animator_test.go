package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Animator
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Give an artifact three +1 power counters. Move it to a flank of its controller's battleline as a creature with versatile for the remainder of the turn.
func TestAnimator(t *testing.T) {
	t.Run("animates a chosen artifact into a 3-power creature for the turn", func(t *testing.T) {
		var animator, target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&animator, Animator),
					ct.Bind(&target, ct.Artifact()),
				),
			},
		})

		h.P1.UseAction(animator)
		h.P1.ClickCard(target) // choose the artifact to animate
		h.P1.ClickOption("left")

		// The chosen artifact — not Animator itself — becomes the creature.
		if got := h.Game().TypeOf(target.ID()); got != engine.Creature {
			t.Fatalf("target should have become a creature, got %v", got)
		}
		if got := h.Game().TypeOf(animator.ID()); got != engine.Artifact {
			t.Errorf("Animator should stay an artifact, got %v", got)
		}
		// Versatile lets the animated creature be used as if in the active house.
		if !h.Game().HasKeyword(target.ID(), engine.Versatile) {
			t.Error("animated creature should gain versatile")
		}
		h.Expect(target).Power(3).At(ct.PlayArea)
	})

	t.Run("reverts to an artifact at end of turn, keeping its counters", func(t *testing.T) {
		var animator, target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&animator, Animator),
					ct.Bind(&target, ct.Artifact()),
				),
			},
		})

		h.P1.UseAction(animator)
		h.P1.ClickCard(target)
		h.P1.ClickOption("left")
		h.Expect(target).Power(3)

		// End of turn reverts the creature to an artifact; its power counters stay.
		h.P1.EndTurn()
		h.P2.EndTurn()
		if got := h.Game().TypeOf(target.ID()); got != engine.Artifact {
			t.Fatalf("target should revert to an artifact at end of turn, got %v", got)
		}
		if err := h.Game().InvariantError(); err != nil {
			t.Fatalf("card conservation broken after revert: %v", err)
		}

		// Re-animating the same artifact stacks three more counters onto the three
		// that persisted, so it enters as a 6-power creature.
		h.P1.ChooseHouse(card.House.Logos)
		h.P1.UseAction(animator)
		h.P1.ClickCard(target)
		h.P1.ClickOption("right")
		h.Expect(target).Power(6).At(ct.PlayArea)
	})
}
