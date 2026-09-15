package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gateway to Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each creature. Gain 3 chains.
func TestGatewayToDis(t *testing.T) {
	t.Run("destroys each creature and gains 3 chains", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				Hand:   ct.Cards(GatewayToDis),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Dis)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Play(GatewayToDis)

		h.Expect(ally).At(ct.Discard)
		h.Expect(foe).At(ct.Discard)
		if got := h.Game().State.Chains[0]; got != 3 {
			t.Fatalf("chains = %d, want 3", got)
		}
	})
}
