package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Stirring Grave
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive a creature from your discard pile.
func TestStirringGrave(t *testing.T) {
	t.Run("archives a creature from your discard pile", func(t *testing.T) {
		var buried ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Dis,
				Hand:    ct.Cards(StirringGrave),
				Discard: ct.Cards(ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Dis)))),
			},
		})

		h.P1.Play(StirringGrave)

		h.Expect(buried).At(ct.Archives)
	})
}
