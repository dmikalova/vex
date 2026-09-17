package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Commandeer
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, each time you play another card, a friendly creature captures 1 Æmber from your opponent.
func TestCommandeer(t *testing.T) {
	var capper, played ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			Hand: ct.Cards(
				Commandeer,
				ct.Bind(&played, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
			),
			InPlay: ct.Cards(
				ct.Bind(&capper, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(4))),
			),
		},
		P2: ct.Side{Amber: 3},
	})

	h.P1.Play(Commandeer)
	// Playing another card fires the reaction; the active player picks a friendly
	// creature to capture with.
	h.P1.Play(played)
	h.P1.ClickCard(capper)

	h.Expect(capper).AmberOn(1)
	h.P2.ExpectAmber(2)
}
