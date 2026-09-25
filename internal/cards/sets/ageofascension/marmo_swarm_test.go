package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Marmo Swarm
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Marmo Swarm gains +1 power for each Æmber in your pool.
func TestMarmoSwarm(t *testing.T) {
	t.Run("gets +1 power for each Æmber in your pool", func(t *testing.T) {
		var marmo ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Amber:  3,
				InPlay: ct.Cards(ct.Bind(&marmo, MarmoSwarm)),
			},
		})

		h.Expect(marmo).Power(5) // 2 + 1 per Æmber, three Æmber
	})

	t.Run("gets no bonus with an empty pool", func(t *testing.T) {
		var marmo ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&marmo, MarmoSwarm)),
			},
		})

		h.Expect(marmo).Power(2)
	})
}
