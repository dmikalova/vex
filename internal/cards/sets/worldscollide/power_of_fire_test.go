package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Power of Fire
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy a friendly Creature -> each player loses Æmber equal to half its power, rounded down. Gain 1 chain.
func TestPowerOfFire(t *testing.T) {
	t.Run(
		"each player loses half the sacrificed creature's power and gains a chain",
		func(t *testing.T) {
			var fodder ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Brobnar,
					Hand:  ct.Cards(PowerOfFire),
					Amber: 5,
					InPlay: ct.Cards(
						ct.Bind(&fodder, ct.Creature(ct.Power(5))),
					),
				},
				P2: ct.Side{Amber: 4},
			})

			h.P1.Play(PowerOfFire)

			h.Expect(fodder).At(ct.Discard)
			h.P1.ExpectAmber(3)
			h.P2.ExpectAmber(2)
			if got := h.Game().State.Chains[0]; got != 1 {
				t.Errorf("chains = %d, want 1", got)
			}
		},
	)
}
