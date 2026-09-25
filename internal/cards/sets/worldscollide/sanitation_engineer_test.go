package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Sanitation Engineer
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Hazardous 1.
//	Reap: Discard a card from your hand.
func TestSanitationEngineer(t *testing.T) {
	t.Run("discards a card from hand when it reaps", func(t *testing.T) {
		var engineer, spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&engineer, SanitationEngineer)),
				Hand:   ct.Cards(ct.Bind(&spare, ct.Creature())),
			},
		})
		engineer.Ready()

		h.P1.Reap(engineer)

		h.Expect(spare).At(ct.Discard)
	})

	t.Run("deals 1 damage to an attacker before combat", func(t *testing.T) {
		var engineer, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(20))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&engineer, SanitationEngineer))},
		})

		h.P1.Fight(foe, engineer)

		h.Expect(foe).Damage(5) // 1 hazardous + 4 retaliation
	})
}
