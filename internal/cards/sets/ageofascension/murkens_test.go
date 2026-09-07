package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Murkens
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Choose one:
//	- Play a random card from your opponent's archives
//	- Play the top card of your opponent's deck.
func TestMurkens(t *testing.T) {
	t.Run("plays the top of the opponent's deck as yours", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Murkens),
			},
			P2: ct.Side{
				Deck: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(Murkens)
		h.P1.ClickOption("top card")

		if !h.Game().InPlay(foe.ID()) {
			t.Fatal("the deck-top creature should be in play")
		}
		if got := h.Game().Controller(foe.ID()); got != 0 {
			t.Errorf("controller = %d, want 0 (played as yours)", got)
		}
		if got := h.Game().Owner(foe.ID()); got != 1 {
			t.Errorf("owner = %d, want 1 (unchanged)", got)
		}
	})

	t.Run("plays a card from the opponent's archives as yours", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Murkens),
			},
			P2: ct.Side{
				Archives: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(Murkens)
		h.P1.ClickOption("random card")

		if !h.Game().InPlay(foe.ID()) {
			t.Fatal("the archived creature should be in play")
		}
		if got := h.Game().Controller(foe.ID()); got != 0 {
			t.Errorf("controller = %d, want 0 (played as yours)", got)
		}
	})

	t.Run("chooses between the two piles when both have cards", func(t *testing.T) {
		var fromDeck, fromArchives ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(Murkens),
			},
			P2: ct.Side{
				Deck:     ct.Cards(ct.Bind(&fromDeck, ct.Creature(ct.Power(3)))),
				Archives: ct.Cards(ct.Bind(&fromArchives, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Play(Murkens)
		h.P1.ClickOption("random card")

		if !h.Game().InPlay(fromArchives.ID()) {
			t.Error("the archived creature should have been played")
		}
		if h.Game().InPlay(fromDeck.ID()) {
			t.Error("the deck-top creature should not have been played")
		}
	})
}
