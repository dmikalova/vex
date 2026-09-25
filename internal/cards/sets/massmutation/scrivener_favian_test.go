package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Scrivener Favian
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	When you resolve a Capture bonus icon, steal 1 Æmber instead.
//	Enhance Capture Capture.
func TestScrivenerFavian(t *testing.T) {
	t.Run("always steals 1 instead of a Capture bonus icon", func(t *testing.T) {
		var bearer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(ScrivenerFavian),
				Hand: ct.Cards(
					ct.Bind(
						&bearer,
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Bonus(card.Bonus.Capture)),
					),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(bearer) // steals instead of capturing, no prompt

		if got := h.P1.Amber(); got != 1 {
			t.Fatalf("P1 Æmber = %d, want 1 (stolen)", got)
		}
		if got := h.P2.Amber(); got != 2 {
			t.Fatalf("P2 Æmber = %d, want 2 (3 - 1 stolen)", got)
		}
	})
}
