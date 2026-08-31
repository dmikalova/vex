package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Annihilation Ritual
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Power
//
//	Each creature gains, "Destroyed: Purge this creature."
func TestAnnihilationRitual(t *testing.T) {
	t.Run("destroyed creatures are purged instead of discarded", func(t *testing.T) {
		var attacker, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					AnnihilationRitual,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(attacker, victim)

		h.Expect(victim).At(ct.Purge)
	})
}
