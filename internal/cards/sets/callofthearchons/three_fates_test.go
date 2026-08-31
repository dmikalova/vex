package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Three Fates
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy the 3 most powerful creatures.
func TestThreeFates(t *testing.T) {
	t.Run("destroys the three most powerful creatures", func(t *testing.T) {
		var big, mid, small, weak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(ThreeFates)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&big, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6))),
				ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				ct.Bind(&small, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				ct.Bind(&weak, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
			)},
		})

		h.P1.Play(ThreeFates)

		h.Expect(big).At(ct.Discard)
		h.Expect(mid).At(ct.Discard)
		h.Expect(small).At(ct.Discard)
		h.Expect(weak).At(ct.PlayArea)
	})
}
