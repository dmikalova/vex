package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Skybooster Squadron
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Martian • Soldier
//
//	Reap: Put Skybooster Squadron into its owner's hand.
func TestSkyboosterSquadron(t *testing.T) {
	t.Run("returns itself to your hand when reaping", func(t *testing.T) {
		var skybooster ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&skybooster, SkyboosterSquadron)),
			},
			P2: ct.Side{},
		})

		h.P1.Reap(skybooster)

		h.Expect(skybooster).At(ct.Hand)
	})
}
