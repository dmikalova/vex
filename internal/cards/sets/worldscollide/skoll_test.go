package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Skoll
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Assault 3.
//	After a Creature is destroyed by Skoll's assault damage, give a friendly Creature a +1 power counter.
func TestSkoll(t *testing.T) {
	t.Run(
		"gives a friendly creature a +1 power counter when its Assault kills",
		func(t *testing.T) {
			var foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Skoll)},
				P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(3))))},
			})

			// Skoll's Assault 3 destroys the power-3 defender before the fight resolves;
			// Skoll is the only friendly creature, so it receives the +1 counter itself.
			h.P1.Fight(Skoll, foe)

			h.Expect(foe).At(ct.Discard)
			h.Expect(Skoll).Power(4)
		},
	)

	t.Run("no counter when the Assault does not kill", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Skoll)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		// Assault 3 only wounds the power-6 defender; the fight then resolves and
		// Skoll (power 3) trades into it, so no power counter is placed.
		h.P1.Fight(Skoll, foe)

		h.Expect(foe).At(ct.Discard) // 3 assault + 3 fight = destroyed by fight damage
		h.Expect(Skoll).At(ct.Discard)
	})
}
