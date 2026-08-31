package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Custom Virus
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Custom Virus. Purge a creature from your hand. Destroy each creature that shares a trait with the purged creature.
func TestCustomVirus(t *testing.T) {
	t.Run("purges a hand creature and destroys creatures sharing its trait", func(t *testing.T) {
		var virus, purged, prey, spared ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&virus, CustomVirus)),
				Hand:   ct.Cards(ct.Bind(&purged, ct.Creature(ct.Traits("Beast")))),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&prey, ct.Creature(ct.Power(5), ct.Traits("Beast"))),
				ct.Bind(&spared, ct.Creature(ct.Power(5), ct.Traits("Robot"))),
			)},
		})

		h.P1.UseAction(virus)

		h.Expect(virus).At(ct.Discard)   // sacrificed
		h.Expect(purged).At(ct.Purge)    // purged from hand
		h.Expect(prey).At(ct.Discard)    // shares Beast -> destroyed
		h.Expect(spared).At(ct.PlayArea) // shares no trait -> survives
	})
}
