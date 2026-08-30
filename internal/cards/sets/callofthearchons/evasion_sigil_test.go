package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Evasion Sigil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Power
//
//	Each creature gains, "Before Fight: Discard the top card of its controller's deck. If it is of the active house, the fight does not occur."
func TestEvasionSigil(t *testing.T) {
	t.Run("cancels the fight when the discarded card is of the active house", func(t *testing.T) {
		var attacker, defender, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					EvasionSigil,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(5))),
				),
				Deck: ct.Cards(ct.Bind(&top, ct.Action(ct.OfHouse(card.House.Shadows)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&defender, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6)))),
			},
		})

		h.P1.Fight(attacker, defender)

		h.Expect(top).At(ct.Discard)
		h.Expect(attacker).Exhausted().Damage(0)
		h.Expect(defender).At(ct.PlayArea).Damage(0)
	})

	t.Run("allows the fight when the discarded card is of a different house", func(t *testing.T) {
		var attacker, defender, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					EvasionSigil,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(5))),
				),
				Deck: ct.Cards(ct.Bind(&top, ct.Action(ct.OfHouse(card.House.Brobnar)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&defender, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6)))),
			},
		})

		h.P1.Fight(attacker, defender)

		h.Expect(top).At(ct.Discard)
		h.Expect(attacker).At(ct.Discard)
		h.Expect(defender).At(ct.PlayArea).Damage(5)
	})
}
