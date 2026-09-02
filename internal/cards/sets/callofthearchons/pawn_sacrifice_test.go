package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Pawn Sacrifice
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy a friendly creature -> deal 3 damage to a creature and deal 3 damage to a different creature.
func TestPawnSacrifice(t *testing.T) {
	t.Run("destroying a friendly creature deals 3 damage each to 2 creatures", func(t *testing.T) {
		var pawn, keeper, foeA, foeB ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(PawnSacrifice),
				InPlay: ct.Cards(
					ct.Bind(&pawn, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(1))),
					ct.Bind(&keeper, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(9))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foeA, ct.Creature(ct.Power(9))),
					ct.Bind(&foeB, ct.Creature(ct.Power(9))),
				),
			},
		})

		h.P1.Play(PawnSacrifice)
		h.P1.ClickCard(pawn)
		h.P1.ClickCard(foeA)
		h.P1.ClickCard(foeB)

		h.Expect(pawn).At(ct.Discard)
		h.Expect(keeper).Damage(0)
		h.Expect(foeA).Damage(3)
		h.Expect(foeB).Damage(3)
	})

	t.Run("with no friendly creature to destroy, no damage is dealt", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(PawnSacrifice),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(PawnSacrifice)

		h.Expect(foe).Damage(0)
	})
}
