package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shorty
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Giant
//
//	Assault 4.
//	Reap: Enrage Shorty.
func TestShorty(t *testing.T) {
	t.Run("enrages itself when it reaps", func(t *testing.T) {
		var shorty ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&shorty, Shorty)),
			},
		})

		h.P1.Reap(shorty)

		h.P1.ExpectAmber(1)
		if !h.Game().Enraged(shorty.ID()) {
			t.Errorf("%s should be enraged", shorty.Name())
		}
	})

	t.Run("deals 4 damage to the attacked enemy before combat", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Shorty),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(20))))},
		})

		h.P1.Fight(Shorty, foe)

		h.Expect(foe).Damage(8)
	})
}
