package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Grim Reminder
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a house. Archive each creature of the chosen house from your discard pile. Gain 1 chain.
func TestGrimReminder(t *testing.T) {
	t.Run("archives each creature of the chosen house from your discard pile", func(t *testing.T) {
		var disCreature, otherHouse, disTactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(GrimReminder),
				Discard: ct.Cards(
					ct.Bind(&disCreature, ct.Creature(ct.OfHouse(card.House.Dis))),
					ct.Bind(&otherHouse, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&disTactic, ct.Tactic(ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.Play(GrimReminder)
		h.P1.ClickOption("Dis")

		h.Expect(disCreature).At(ct.Archives)
		h.Expect(otherHouse).At(ct.Discard)
		h.Expect(disTactic).At(ct.Discard)
		if got := h.Game().State.Chains[0]; got != 1 {
			t.Errorf("chains = %d, want 1", got)
		}
	})

	t.Run("gains a chain even when no creature matches", func(t *testing.T) {
		var otherHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(GrimReminder),
				Discard: ct.Cards(
					ct.Bind(&otherHouse, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.Play(GrimReminder)
		h.P1.ClickOption("Dis")

		h.Expect(otherHouse).At(ct.Discard)
		if got := h.Game().State.Chains[0]; got != 1 {
			t.Errorf("chains = %d, want 1", got)
		}
	})
}
