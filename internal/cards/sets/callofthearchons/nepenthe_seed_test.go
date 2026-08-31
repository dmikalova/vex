package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nepenthe Seed
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Nepenthe Seed, and put a card from your discard pile into your hand.
func TestNepentheSeed(t *testing.T) {
	t.Run("sacrifices itself and returns a card from the discard pile to hand", func(t *testing.T) {
		var ghost ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(NepentheSeed),
				Discard: ct.Cards(
					ct.Bind(&ghost, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(4))),
				),
			},
		})

		h.P1.UseAction(NepentheSeed)
		h.P1.ClickCard(ghost) // Nepenthe Seed is also in the discard now; retrieve the ghost

		h.Expect(NepentheSeed).At(ct.Discard) // sacrificed itself
		h.Expect(ghost).At(ct.Hand)           // returned from discard
	})
}
