package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Arise!
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a house - put each creature of the chosen house from your discard pile into your hand, and gain 1 chain.
func TestArise(t *testing.T) {
	var dis1, dis2, sanc ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(Arise),
			Discard: ct.Cards(
				ct.Bind(&dis1, ct.Creature(ct.OfHouse(card.House.Dis))),
				ct.Bind(&dis2, ct.Creature(ct.OfHouse(card.House.Dis))),
				ct.Bind(&sanc, ct.Creature(ct.OfHouse(card.House.Sanctum))),
			),
		},
	})

	h.P1.Play(Arise)
	h.P1.ClickOption("Dis")

	h.Expect(dis1).At(ct.Hand)
	h.Expect(dis2).At(ct.Hand)
	h.Expect(sanc).At(ct.Discard)
	if got := h.Game().State.Chains[0]; got != 1 {
		t.Errorf("chains = %d, want 1", got)
	}
}
