package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// EMP Blast
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Stun each Mars creature and each Robot trait creature, and destroy each artifact.
func TestEMPBlast(t *testing.T) {
	t.Run("stuns each Mars and Robot creature and destroys each artifact", func(t *testing.T) {
		var marsGuy, robot, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(EMPBlast)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&marsGuy, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&robot, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3), ct.Traits("Robot"))),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
					Cannon,
				),
			},
		})

		h.P1.Play(EMPBlast)

		h.Expect(marsGuy).Stunned(true)
		h.Expect(robot).Stunned(true)
		h.Expect(other).Stunned(false)
		h.Expect(Cannon).At(ct.Discard) // each artifact is destroyed
	})
}
