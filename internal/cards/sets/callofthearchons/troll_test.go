package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Troll
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  8
//	Traits: Giant
//
//	Reap: Heal 3 damage from Troll.
func TestTroll(t *testing.T) {
	t.Run("heals 3 damage from itself when it reaps", func(t *testing.T) {
		var troll ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(ct.Bind(&troll, Troll))},
		})

		troll.Damaged(3)
		h.P1.Reap(Troll)

		h.Expect(troll).Damage(0)
	})
}
