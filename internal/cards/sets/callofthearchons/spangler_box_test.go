package callofthearchons

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Spangler Box
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Graft a creature from play, and your opponent gains control of Spangler Box.
//	Destroyed: Put each card under Spangler Box into play under its owner's control.
func TestSpanglerBox(t *testing.T) {
	t.Run("action grafts a creature and hands the box to the opponent", func(t *testing.T) {
		var box, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&box, SpanglerBox),
					ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
			P2: ct.Side{House: card.House.Brobnar},
		})

		h.P1.UseAction(box)

		h.Expect(box).At(ct.PlayArea)
		h.Expect(victim).At(ct.Under)
		if got := h.Game().Controller(box.ID()); got != 1 {
			t.Fatalf("box controller = %d, want 1 (the opponent)", got)
		}
	})

	t.Run("destroying the box returns grafted creatures to their owner", func(t *testing.T) {
		var box, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&box, SpanglerBox),
					ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
			P2: ct.Side{House: card.House.Brobnar},
		})

		h.P1.UseAction(box)
		h.Expect(victim).At(ct.Under)

		g := h.Game()
		g.DestroyEach(g.Controller(box.ID()), []engine.LocalID{box.ID()})

		h.Expect(box).At(ct.Discard)
		h.Expect(victim).At(ct.PlayArea)
		if !slices.Contains(g.Battleline(0), victim.ID()) {
			t.Fatalf(
				"victim is in battlelines %v/%v, want P1's",
				g.Battleline(0),
				g.Battleline(1),
			)
		}
	})
}
