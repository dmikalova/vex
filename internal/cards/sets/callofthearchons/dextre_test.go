package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Dextre
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Play: Dextre captures 1 Æmber from your opponent.
//	Destroyed: Put Dextre on top of its owner's deck.
func TestDextre(t *testing.T) {
	t.Run("captures 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(Dextre),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(Dextre)

		h.Expect(Dextre).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("returns to the top of its owner's deck when destroyed", func(t *testing.T) {
		var dextre ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&dextre, Dextre)),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{dextre.ID()})

		h.Expect(dextre).At(ct.Deck) // top of deck, not the discard
	})
}
