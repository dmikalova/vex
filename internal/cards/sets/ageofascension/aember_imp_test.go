package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Aember Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	After a creature reaps, stun it.
func TestAemberImp(t *testing.T) {
	t.Run("stuns the reaper, even itself", func(t *testing.T) {
		var imp ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&imp, AemberImp)),
			},
		})

		h.P1.Reap(imp)

		h.Expect(imp).Stunned(true)
	})
}
