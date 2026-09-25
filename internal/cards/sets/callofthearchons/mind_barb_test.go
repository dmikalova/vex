package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Mind Barb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Discard a card from your hand. Your opponent discards a random card from their hand.
func TestMindBarb(t *testing.T) {
	t.Run("discards a chosen card and the opponent discards a random one", func(t *testing.T) {
		var mine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(MindBarb, ct.Bind(&mine, ct.Creature())),
			},
			P2: ct.Side{Hand: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Play(MindBarb)

		h.Expect(mine).At(ct.Discard) // the only other card in hand, chosen to discard
		if got := h.Game().State.Hand[1].Count; got != 1 {
			t.Fatalf("opponent hand = %d, want 1", got)
		}
		if got := h.Game().State.Discard[1].Count; got != 1 {
			t.Fatalf("opponent discard = %d, want 1", got)
		}
	})
}
