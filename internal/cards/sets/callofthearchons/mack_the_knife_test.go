package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mack the Knife
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive, Versatile.
//	Action: Deal 1 damage to a creature. If this damage destroys that creature, gain 1 Æmber.
func TestMackTheKnife(t *testing.T) {
	t.Run("gains 1 Æmber when its damage destroys the creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(MackTheKnife)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))),
				),
			},
		})

		h.P1.UseAction(MackTheKnife)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})
}
