package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Seeker Needle
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Deal 1 damage to a creature. If this damage destroys that creature, gain 1 Æmber.
func TestSeekerNeedle(t *testing.T) {
	t.Run("gains 1 Æmber when its damage destroys the creature", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(SeekerNeedle)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(1))))},
		})

		h.P1.UseAction(SeekerNeedle)

		h.Expect(foe).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})

	t.Run("gains nothing when the creature survives", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(SeekerNeedle)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))))},
		})

		h.P1.UseAction(SeekerNeedle)

		h.Expect(foe).At(ct.PlayArea).Damage(1)
		h.P1.ExpectAmber(0)
	})
}
