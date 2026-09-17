package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ascendant Hester
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  8
//	Traits: Knight • Spirit
//
//	Each other friendly creature gains +2 armor for each Æmber on it.
//	Play/Fight: Each friendly creature captures 1 Æmber from your opponent.
func TestAscendantHester(t *testing.T) {
	var hester, ally ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Sanctum,
			Hand:   ct.Cards(ct.Bind(&hester, AscendantHester), card.GiganticArt(AscendantHester)),
			InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(4)))),
		},
		P2: ct.Side{Amber: 3},
	})

	h.P1.Play(hester)

	// Each friendly creature captured 1 Æmber from the opponent.
	h.P2.ExpectAmber(1)
	h.Expect(hester).AmberOn(1)
	h.Expect(ally).AmberOn(1)
	// The ally now carries 1 Æmber, so Hester grants it +2 armor for that Æmber.
	h.Expect(ally).Armor(2)
}
