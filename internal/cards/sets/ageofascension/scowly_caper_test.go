package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Scowly Caper
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Skirmish, Treachery, Versatile.
//	At the end of your turn, destroy a neighboring creature.
func TestScowlyCaper(t *testing.T) {
	t.Run("enters play under your opponent's control", func(t *testing.T) {
		var scowly ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ct.Bind(&scowly, ScowlyCaper)),
			},
		})

		h.P1.Play(ScowlyCaper)

		if !h.Game().InPlay(scowly.ID()) {
			t.Fatal("Scowly Caper should be in play")
		}
		if got := h.Game().Controller(scowly.ID()); got != 1 {
			t.Errorf("controller = %d, want 1 (your opponent)", got)
		}
		if got := h.Game().Owner(scowly.ID()); got != 0 {
			t.Errorf("owner = %d, want 0 (unchanged)", got)
		}
	})

	t.Run("destroys a neighbor at the end of its controller's turn", func(t *testing.T) {
		var victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(ScowlyCaper),
			},
			P2: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&victim,
					ct.Creature(ct.Power(3), ct.OfHouse(card.House.Shadows)))),
			},
		})

		// Scowly Caper enters under P2's control, on the flank beside P2's creature.
		h.P1.Play(ScowlyCaper)

		// The destroy fires at the end of the controller's (P2's) turn.
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Shadows)
		h.P2.EndTurn()

		if h.Game().InPlay(victim.ID()) {
			t.Error("the neighbor should have been destroyed")
		}
	})
}
