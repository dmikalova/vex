package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Auto-Encoder
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	After you discard a card from your hand, archive the top card of your deck.
func TestAutoEncoder(t *testing.T) {
	t.Run(
		"archives the top card of your deck after a card is discarded from hand",
		func(t *testing.T) {
			var discarded, deckTop ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Logos,
					InPlay: ct.Cards(AutoEncoder),
					Hand:   ct.Cards(ct.Bind(&discarded, ct.Tactic())),
					Deck: ct.Cards(
						ct.Bind(&deckTop, ct.Tactic()),
						ct.Tactic(),
					),
				},
			})

			h.P1.Discard(discarded)

			h.Expect(deckTop).At(ct.Archives)
		},
	)
}
