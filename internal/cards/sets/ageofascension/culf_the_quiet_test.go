package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Culf the Quiet
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant
//
//	Elusive.
func TestCulfTheQuiet(t *testing.T) {
	t.Run("elusive: takes no damage the first time it is attacked", func(t *testing.T) {
		var attacker ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&attacker, ct.Creature(ct.Power(5)))),
			},
			P2: ct.Side{InPlay: ct.Cards(CulfTheQuiet)},
		})

		h.P1.Fight(attacker, CulfTheQuiet)

		h.Expect(CulfTheQuiet).Damage(0)
	})
}
