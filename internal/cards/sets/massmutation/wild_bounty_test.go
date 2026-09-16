package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Wild Bounty
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: The next time you play a card this turn, resolve each of its bonus icons an additional time.
//	Enhance Æmber Æmber.
func TestWildBounty(t *testing.T) {
	t.Run("next played card resolves each bonus icon an additional time", func(t *testing.T) {
		var bounty, budded ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand: ct.Cards(
					ct.Bind(&bounty, WildBounty),
					ct.Bind(&budded, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Bonus(card.Bonus.Aember))),
				),
			},
		})

		h.P1.Play(bounty)
		before := h.P1.Amber()
		h.P1.Play(budded)

		// The creature's single Æmber icon resolves twice: base once, boost once.
		if got := h.P1.Amber() - before; got != 2 {
			t.Fatalf("aember gained = %d, want 2", got)
		}
	})

	t.Run("boost fires for one card only", func(t *testing.T) {
		var bounty, first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand: ct.Cards(
					ct.Bind(&bounty, WildBounty),
					ct.Bind(&first, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Bonus(card.Bonus.Aember))),
					ct.Bind(&second, ct.Creature(
						ct.OfHouse(card.House.Untamed), ct.Bonus(card.Bonus.Aember))),
				),
			},
		})

		h.P1.Play(bounty)
		h.P1.Play(first)
		before := h.P1.Amber()
		h.P1.Play(second)

		// The boost was spent on the first card; the second gains its icon once.
		if got := h.P1.Amber() - before; got != 1 {
			t.Fatalf("second card aember = %d, want 1", got)
		}
	})
}
