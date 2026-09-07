package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Doctor Driscoll
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Scientist
//
//	Elusive.
//	Action: Heal 2 damage from a creature. For each damage healed this way, gain 1 Æmber.
func TestDoctorDriscoll(t *testing.T) {
	t.Run("heals 2 damage from a creature and gains 1 Æmber per damage healed", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					DoctorDriscoll,
					ct.Bind(&target, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(6))),
				),
			},
		})
		target.Damaged(3)

		h.P1.UseAction(DoctorDriscoll)
		h.P1.ClickCard(target)

		h.Expect(target).Damage(1)
		h.P1.ExpectAmber(2)
	})
}
