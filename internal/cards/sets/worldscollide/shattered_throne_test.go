package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shattered Throne
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	After a creature is used to fight, it captures 1 Æmber from its opponent.
func TestShatteredThrone(t *testing.T) {
	t.Run("the fighting creature captures 1 Æmber after it fights", func(t *testing.T) {
		var attacker ct.Card
		var defender ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ShatteredThrone,
					ct.Bind(&attacker, ct.Creature(ct.Power(5))),
				),
			},
			P2: ct.Side{
				Amber:  3,
				InPlay: ct.Cards(ct.Bind(&defender, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Fight(attacker, defender)

		h.Expect(attacker).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
