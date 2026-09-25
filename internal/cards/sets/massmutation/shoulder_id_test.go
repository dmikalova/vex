package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Shoulder Id
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Specter
//
//	Taunt.
//	When Shoulder Id would deal damage, steal 1 Æmber instead.
//	Shoulder Id cannot fight.
func TestShoulderId(t *testing.T) {
	var shoulder, attacker ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Shadows,
			Amber: 3,
			InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(
				ct.OfHouse(card.House.Shadows), ct.Power(5)))),
		},
		P2: ct.Side{
			InPlay: ct.Cards(ct.Bind(&shoulder, ShoulderID)),
		},
	})

	h.P1.Fight(attacker, shoulder)

	// Shoulder Id retaliates by stealing, not by dealing damage.
	h.Expect(attacker).Damage(0)
	h.Expect(shoulder).Damage(5)
	h.P2.ExpectAmber(1)
	h.P1.ExpectAmber(2)

	// It cannot fight.
	h.P2.ExpectCannotUseTo(shoulder, card.UseKind.Fight)
}
