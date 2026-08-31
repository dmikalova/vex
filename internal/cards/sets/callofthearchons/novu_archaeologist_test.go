package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Novu Archaeologist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Action: Archive a card from your discard pile.
func TestNovuArchaeologist(t *testing.T) {
	t.Run("archives a chosen card from the discard pile", func(t *testing.T) {
		var buried, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(NovuArchaeologist),
				Discard: ct.Cards(
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Logos))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
			},
		})

		h.P1.UseAction(NovuArchaeologist)
		h.P1.ClickCard(buried)

		h.Expect(buried).At(ct.Archives)
		h.Expect(other).At(ct.Discard)
	})
}
