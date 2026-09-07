package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Eldest Bear
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Beast • Leader • Witch
//
//	Assault 3.
//	While Eldest Bear is in the center of your battleline, it gains, "Before Fight: Gain 2 Æmber."
func TestEldestBear(t *testing.T) {
	t.Run("centered, gains 2 Æmber before it fights", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(
				ct.Creature(),
				EldestBear,
				ct.Creature(),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(EldestBear, enemy)

		h.P1.ExpectAmber(2)
	})

	t.Run("off-center, no Before Fight Æmber", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(
				EldestBear,
				ct.Creature(),
			)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(EldestBear, enemy)

		h.P1.ExpectAmber(0)
	})
}
