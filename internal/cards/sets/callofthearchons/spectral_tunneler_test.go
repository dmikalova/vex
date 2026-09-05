package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Spectral Tunneler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Choose a creature - for the remainder of the turn, it is considered a flank creature, and it gains, "Reap: Draw a card."
func TestSpectralTunneler(t *testing.T) {
	// The chosen creature keeps its normal reap (1 Æmber) and, thanks to the
	// granted "Reap: Draw a card", also draws.
	t.Run("the chosen creature draws when it reaps", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					SpectralTunneler,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Deck: ct.Cards(ct.Creature(), ct.Creature()),
			},
		})

		before := len(h.Game().Hand(0))
		h.P1.UseAction(SpectralTunneler)
		h.P1.Reap(ally)

		if got := len(h.Game().Hand(0)); got != before+1 {
			t.Errorf("hand after reap = %d, want %d (the granted Reap ability drew)", got, before+1)
		}
		h.P1.ExpectAmber(1) // the reap itself still gains 1 Æmber
	})

	// The grant is scoped to the one chosen creature, so a different creature's
	// reap does not draw.
	t.Run("only the chosen creature draws on reap", func(t *testing.T) {
		var chosen, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					SpectralTunneler,
					ct.Bind(&chosen, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Deck: ct.Cards(ct.Creature(), ct.Creature()),
			},
		})

		h.P1.UseAction(SpectralTunneler)
		h.P1.ClickCard(chosen) // two creatures on the board, so the choice is offered

		before := len(h.Game().Hand(0))
		h.P1.Reap(other)
		if got := len(h.Game().Hand(0)); got != before {
			t.Errorf("hand after a different creature reaps = %d, want %d (no draw)", got, before)
		}
	})
}
