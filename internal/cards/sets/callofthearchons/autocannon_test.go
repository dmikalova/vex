package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Autocannon
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	After a creature enters play, deal 1 damage to it.
func TestAutocannon(t *testing.T) {
	t.Run("deals 1 damage to a creature as it enters play", func(t *testing.T) {
		var newcomer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(Autocannon),
				Hand: ct.Cards(
					ct.Bind(&newcomer, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(newcomer) // Autocannon zaps the newcomer for 1

		h.Expect(newcomer).Damage(1)
	})
}
