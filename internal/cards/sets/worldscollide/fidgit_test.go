package worldscollide

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Fidgit
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Faerie • Thief
//
//	Elusive.
//	Reap: Discard a random card from your opponent's archives or the top card of their deck. If that card is a tactic, play it as if it were yours.
func TestFidgit(t *testing.T) {
	t.Run("discards and plays a Tactic from the opponent's archives", func(t *testing.T) {
		var tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(Fidgit),
			},
			P2: ct.Side{
				Archives: ct.Cards(ct.Bind(&tactic, ct.Tactic(ct.AemberBonus(1)))),
			},
		})

		before := h.Game().Aember(0)
		h.P1.Reap(Fidgit)
		h.P1.ClickOption("archives")

		if slices.Contains(h.Game().Archives(1), tactic.ID()) {
			t.Error("the tactic should have left the opponent's archives")
		}
		if got := h.Game().Aember(0) - before; got != 2 {
			t.Errorf("Æmber gained = %d, want 2 (1 reap + 1 tactic pip played as yours)", got)
		}
	})

	t.Run("discards but does not play a non-Tactic from the deck", func(t *testing.T) {
		var creature ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(Fidgit),
			},
			P2: ct.Side{
				Deck: ct.Cards(ct.Bind(&creature, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Reap(Fidgit)
		h.P1.ClickOption("top card")

		if h.Game().InPlay(creature.ID()) {
			t.Error("a non-Tactic should be discarded, not played")
		}
		if !slices.Contains(h.Game().Discard(1), creature.ID()) {
			t.Error("the discarded card should be in the opponent's discard pile")
		}
	})
}
