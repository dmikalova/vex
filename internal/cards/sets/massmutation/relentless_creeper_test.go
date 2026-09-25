package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Relentless Creeper
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	After you choose Dis as your active house, you may put Relentless Creeper from your discard pile into your hand.
func TestRelentlessCreeper(t *testing.T) {
	// A stocked deck keeps the draw phase from recycling the discard pile, so the
	// Creeper stays in the discard until its own ability returns it.
	setup := func() ct.Setup {
		return ct.Setup{
			P1: ct.Side{
				House:   card.House.Dis,
				Discard: ct.Cards(RelentlessCreeper),
				Deck:    ct.DeckOf(card.House.Dis, 10),
			},
		}
	}

	// cycle advances from the opening turn (where the house is chosen before cards
	// are placed) round to P1's next turn, then chooses house for P1.
	cycle := func(h *ct.Harness, house engine.House) {
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(house)
	}

	t.Run("returns itself from the discard pile to hand when Dis is chosen", func(t *testing.T) {
		h := ct.Play(t, setup())
		cycle(h, card.House.Dis)
		h.P1.ClickOption("Yes")

		h.Expect(RelentlessCreeper).At(ct.Hand)
	})

	t.Run("stays in the discard pile when the return is declined", func(t *testing.T) {
		h := ct.Play(t, setup())
		cycle(h, card.House.Dis)
		h.P1.ClickOption("No")

		h.Expect(RelentlessCreeper).At(ct.Discard)
	})

	t.Run("does nothing when a different house is chosen", func(t *testing.T) {
		h := ct.Play(t, setup())
		cycle(h, card.House.Untamed)

		h.Expect(RelentlessCreeper).At(ct.Discard)
	})
}
