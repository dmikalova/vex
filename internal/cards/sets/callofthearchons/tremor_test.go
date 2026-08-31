package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tremor
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Stun a creature and each of its neighbors.
func TestTremor(t *testing.T) {
	t.Run("stuns a chosen creature and each of its neighbors", func(t *testing.T) {
		var left, mid, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Tremor)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				ct.Bind(&mid, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		h.P1.Play(Tremor)
		h.P1.ClickCard(mid)

		h.Expect(left).Stunned(true)
		h.Expect(mid).Stunned(true)
		h.Expect(right).Stunned(true)
	})
}
