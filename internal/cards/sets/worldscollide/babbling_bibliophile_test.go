package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Babbling Bibliophile
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Cyborg • Scientist
//
//	Reap: Draw 2 cards.
func TestBabblingBibliophile(t *testing.T) {
	t.Run("draws 2 cards when it reaps", func(t *testing.T) {
		var biblio ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&biblio, BabblingBibliophile)),
				Deck:   ct.Cards(ct.Creature(), ct.Creature()),
			},
		})
		biblio.Ready()

		h.P1.Reap(biblio)

		h.P1.ExpectAmber(1)
		if got := len(h.Game().Hand(0)); got != 2 {
			t.Errorf("hand = %d cards, want 2", got)
		}
	})
}
