package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Effervescent Principle
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Each player loses half of their Æmber, rounded down, and gain 1 chain.
func TestEffervescentPrinciple(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(EffervescentPrinciple), Amber: 5},
		P2: ct.Side{Amber: 4},
	})

	h.P1.Play(EffervescentPrinciple)

	h.P1.ExpectAmber(3) // loses 2 (floor 5/2)
	h.P2.ExpectAmber(2) // loses 2 (floor 4/2)
	if got := h.Game().State.Chains[0]; got != 1 {
		t.Errorf("chains = %d, want 1", got)
	}
}
