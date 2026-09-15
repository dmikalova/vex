package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Alaka
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Giant
//
//	If you have used a creature to fight this turn, Alaka enters play ready.
func TestAlaka(t *testing.T) {
	t.Run("enters play ready once you have fought this turn", func(t *testing.T) {
		var alaka, attacker, enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&alaka, Alaka)),
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(3))))},
		})

		h.P1.Fight(attacker, enemy)
		h.P1.Play(Alaka)

		if alaka.Exhausted() {
			t.Error("Alaka should enter play ready after a fight")
		}
	})

	t.Run("enters play exhausted when you have not fought", func(t *testing.T) {
		var alaka ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(ct.Bind(&alaka, Alaka))},
		})

		h.P1.Play(Alaka)

		if !alaka.Exhausted() {
			t.Error("Alaka should enter play exhausted with no fight")
		}
	})
}
