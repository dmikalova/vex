package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hold the Line
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For each creature your opponent controls in excess of you, draw a card.
func TestHoldTheLine(t *testing.T) {
	t.Run("draws nothing when not outnumbered", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				Hand:   ct.Cards(HoldTheLine),
				InPlay: ct.Cards(ct.Creature()),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature())},
		})

		h.P1.Play(HoldTheLine)

		h.Expect(top).At(ct.Deck)
	})

	t.Run("draws one when outnumbered by one", func(t *testing.T) {
		var first, second ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				Hand:   ct.Cards(HoldTheLine),
				InPlay: ct.Cards(ct.Creature()),
				Deck:   ct.Cards(ct.Bind(&first, ct.Creature()), ct.Bind(&second, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Play(HoldTheLine)

		h.Expect(first).At(ct.Hand)
		h.Expect(second).At(ct.Deck)
	})

	t.Run("draws the full difference when outnumbered by two", func(t *testing.T) {
		var first, second, third ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand:  ct.Cards(HoldTheLine),
				Deck: ct.Cards(
					ct.Bind(&first, ct.Creature()),
					ct.Bind(&second, ct.Creature()),
					ct.Bind(&third, ct.Creature()),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(), ct.Creature())},
		})

		h.P1.Play(HoldTheLine)

		h.Expect(first).At(ct.Hand)
		h.Expect(second).At(ct.Hand)
		h.Expect(third).At(ct.Deck)
	})
}
