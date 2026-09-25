package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Break-key
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If your opponent has more forged keys than you, unforge one of your opponent's keys, and your opponent gains 6 Æmber.
func TestBreakkey(t *testing.T) {
	t.Run("unforges an opponent key and gives them 6 when they lead on keys", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Breakkey),
			},
			P2: ct.Side{ForgedKeys: []engine.KeyColor{card.KeyColor.Red}},
		})

		h.P1.Play(Breakkey)
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(6)
	})

	t.Run("does nothing when the opponent is not ahead on keys", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Breakkey),
			},
		})

		h.P1.Play(Breakkey)
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(0)
	})
}
