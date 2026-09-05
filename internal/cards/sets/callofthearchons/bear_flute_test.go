package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bear Flute
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Fully heal an Ancient Bear. If there are no Ancient Bears in play, search your deck and discard pile and put each Ancient Bear from them into your hand -> shuffle your discard pile into your deck.
func TestBearFlute(t *testing.T) {
	t.Run("heals a bear that is in play", func(t *testing.T) {
		var flute, bear ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&flute, BearFlute),
					ct.Bind(&bear, AncientBear),
				),
				Deck: ct.Cards(AncientBear),
			},
		})
		bear.Damaged(3)

		h.P1.UseAction(flute)

		h.Expect(bear).Damage(0)
		h.Expect(bear).At(ct.PlayArea)
	})

	t.Run("with no bear in play, gathers every copy and reshuffles", func(t *testing.T) {
		var flute, inDeck, inDiscard ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Untamed,
				InPlay:  ct.Cards(ct.Bind(&flute, BearFlute)),
				Deck:    ct.Cards(ct.Bind(&inDeck, AncientBear), ct.Creature()),
				Discard: ct.Cards(ct.Bind(&inDiscard, AncientBear), ct.Creature()),
			},
		})

		h.P1.UseAction(flute)

		if inDeck.Location() != ct.Hand || inDiscard.Location() != ct.Hand {
			t.Errorf("bears are at %v and %v, want both in hand",
				inDeck.Location(), inDiscard.Location())
		}
		if n := len(h.Game().Discard(0)); n != 0 {
			t.Errorf("discard pile holds %d cards, want 0 after the reshuffle", n)
		}
	})
}
