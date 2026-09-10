package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chain Gang
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	After you play Subtle Chain, ready Chain Gang.
//	Action: Steal 1 Æmber. Shuffle a Subtle Chain from your discard pile into your deck.
func TestChainGang(t *testing.T) {
	t.Run("playing Subtle Chain readies Chain Gang", func(t *testing.T) {
		var gang ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&gang, ChainGang)),
				Hand:   ct.Cards(SubtleChain),
			},
			P2: ct.Side{Hand: ct.Cards(ct.Creature())},
		})
		gang.Exhaust()

		h.P1.Play(SubtleChain)

		h.Expect(gang).Ready()
	})

	t.Run("a different play does not ready Chain Gang", func(t *testing.T) {
		var gang ct.Card
		filler := ct.Creature(ct.OfHouse(card.House.Shadows))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&gang, ChainGang)),
				Hand:   ct.Cards(filler),
			},
		})
		gang.Exhaust()

		h.P1.Play(filler)

		h.Expect(gang).Exhausted()
	})

	t.Run("Action steals 1 and shuffles a Subtle Chain into the deck", func(t *testing.T) {
		var gang ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Shadows,
				InPlay:  ct.Cards(ct.Bind(&gang, ChainGang)),
				Discard: ct.Cards(SubtleChain),
			},
			P2: ct.Side{Amber: 3},
		})
		deckBefore := h.Game().State.Deck[0].Count

		h.P1.UseAction(gang)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
		if got := h.Game().State.Discard[0].Count; got != 0 {
			t.Fatalf("discard = %d, want 0 (Subtle Chain shuffled away)", got)
		}
		if got := h.Game().State.Deck[0].Count; got != deckBefore+1 {
			t.Fatalf("deck = %d, want %d", got, deckBefore+1)
		}
	})
}
