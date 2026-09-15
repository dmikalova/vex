package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bulleteye
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: Destroy a flank creature.
func TestBulleteye(t *testing.T) {
	t.Run("Reap destroys a chosen flank creature, not a middle one", func(t *testing.T) {
		var left, middle ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(Bulleteye),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&middle, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
					ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5)),
				),
			},
		})

		h.P1.Reap(Bulleteye)
		h.P1.ClickCard(left)

		h.Expect(left).At(ct.Discard)
		h.Expect(middle).At(ct.PlayArea)
	})
}
