package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Foggify
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: During your opponent's next turn, after an enemy Creature is used to fight, stun it.
func TestFoggify(t *testing.T) {
	t.Run("arms the opponent's next turn without restricting the caster", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Logos, Hand: ct.Cards(Foggify)},
		})

		h.P1.Play(Foggify)

		if !h.Game().State.StunFighterNext[1].Value {
			t.Error("Foggify should arm the opponent's next turn")
		}
		if h.Game().State.StunFighterNext[0].Value {
			t.Error("Foggify should not arm the caster")
		}
	})

	t.Run("stuns each enemy creature the opponent uses to fight next turn", func(t *testing.T) {
		var attacker, defender ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				Hand:   ct.Cards(Foggify),
				InPlay: ct.Cards(ct.Bind(&defender, ct.Creature(ct.Power(2)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
				),
			},
		})

		h.P1.Play(Foggify)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Fight(attacker, defender)

		h.Expect(attacker).Stunned(true)
	})
}
