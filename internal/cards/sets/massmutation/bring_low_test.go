package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bring Low
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Capture all but 5 Æmber from your opponent, distributed among any number of friendly creatures.
//	Enhance Capture.
func TestBringLow(t *testing.T) {
	var bring, a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			Hand:  ct.Cards(ct.Bind(&bring, BringLow)),
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.Power(4))),
				ct.Bind(&b, ct.Creature(ct.Power(4))),
			),
		},
		P2: ct.Side{Amber: 8},
	})

	h.P1.Play(bring)
	// All but 5 of the opponent's 8 Æmber is 3, distributed one at a time.
	h.P1.ClickCard(a)
	h.P1.ClickCard(a)
	h.P1.ClickCard(b)

	h.P2.ExpectAmber(5)
	h.Expect(a).AmberOn(2)
	h.Expect(b).AmberOn(1)
}
