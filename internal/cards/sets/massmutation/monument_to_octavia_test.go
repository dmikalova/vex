package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Monument to Octavia
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: If Cornicen Octavia is in your discard pile, a friendly creature captures 2 Æmber from your opponent. Otherwise, a friendly creature captures 1 Æmber from your opponent.
func TestMonumentToOctavia(t *testing.T) {
	t.Run("a friendly creature captures 1 Æmber", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					MonumentToOctavia,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(MonumentToOctavia)

		h.Expect(ally).AmberOn(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("captures 2 when Cornicen Octavia is in your discard pile", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					MonumentToOctavia,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
				),
				Discard: ct.Cards(CornicenOctavia),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(MonumentToOctavia)

		h.Expect(ally).AmberOn(2)
		h.P2.ExpectAmber(1)
	})
}
