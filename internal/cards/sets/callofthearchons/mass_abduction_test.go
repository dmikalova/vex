package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mass Abduction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Put up to 3 enemy damaged creatures into your archives.
func TestMassAbduction(t *testing.T) {
	t.Run("abducts the damaged enemy creatures the controller chooses", func(t *testing.T) {
		var abduction, hurt, alsoHurt, healthy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				Hand:  ct.Cards(ct.Bind(&abduction, MassAbduction)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&hurt, ct.Creature(ct.Power(5))),
					ct.Bind(&alsoHurt, ct.Creature(ct.Power(6))),
					ct.Bind(&healthy, ct.Creature(ct.Power(7))),
				),
			},
		})

		hurt.Damaged(2)
		alsoHurt.Damaged(1)

		h.P1.Play(abduction)
		h.P1.ClickCard(hurt)
		h.P1.ClickCard(alsoHurt)

		h.Expect(hurt).At(ct.Archives)
		h.Expect(alsoHurt).At(ct.Archives)
		h.Expect(healthy).At(ct.PlayArea)
	})
}
