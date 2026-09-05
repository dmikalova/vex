package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Howling Pit
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	During their "draw cards" step, each player refills their hand to 1 additional card.
func TestTheHowlingPit(t *testing.T) {
	t.Run("refills each player's hand to one additional card", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(TheHowlingPit),
				Deck:   ct.DeckOf(card.House.Logos, 10),
			},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 7 {
			t.Errorf("hand after draw = %d, want 7", got)
		}
	})
}
