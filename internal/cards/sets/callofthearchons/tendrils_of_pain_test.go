package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tendrils of Pain
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 1 damage to each creature, and if your opponent forged a key on their previous turn, deal 3 damage to each creature.
func TestTendrilsOfPain(t *testing.T) {
	t.Run("deals only 1 damage when the opponent forged nothing", func(t *testing.T) {
		var tendrils, beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&tendrils, TendrilsOfPain)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&beast, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Play(tendrils)
		h.Expect(beast).Damage(1)
	})

	t.Run("deals 4 damage when the opponent forged a key last turn", func(t *testing.T) {
		var tendrils, beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&tendrils, TendrilsOfPain)),
			},
			P2: ct.Side{
				Amber: 6,
				InPlay: ct.Cards(
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(9))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.ExpectKeys(1)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Dis)
		h.P1.Play(tendrils)

		h.Expect(beast).Damage(4)
	})
}
