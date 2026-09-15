package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sagittarii's Gaze
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Exalt a damaged creature.
//	Enhance Damage.
func TestSagittariisGaze(t *testing.T) {
	var wounded ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(SagittariisGaze),
			InPlay: ct.Cards(
				ct.Bind(&wounded, ct.Creature(ct.Power(5))),
			),
		},
	})
	wounded.Damaged(2)

	h.P1.Play(SagittariisGaze) // only the damaged creature can be exalted

	h.Expect(wounded).AmberOn(1)
}
